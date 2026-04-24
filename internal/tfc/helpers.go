package tfc

import (
    "bufio"
    "encoding/json"
    "fmt"
    "github.com/fatih/color"
    "github.com/hashicorp/go-retryablehttp"
    "github.com/hashicorp/go-tfe"
    "github.com/hashicorp/jsonapi"
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"
)

type logMessage struct {
    Level      string      `json:"@level"`
    Message    string      `json:"@message"`
    Module     string      `json:"@module"`
    Timestamp  string      `json:"@timestamp"`
    Change     interface{} `json:"change"`
    Diagnostic interface{} `json:"diagnostic"`
    Snippet    interface{} `json:"snippet"`
    Type       string      `json:"type"`
    Terraform  string      `json:"terraform"`
    UI         string      `json:"ui"`
}

// LogRead https://developer.hashicorp.com/terraform/internals/machine-readable-ui
func renderLog(logs io.Reader, err error) {
    if err != nil {
        log.Fatal("Failed to retrieve logs", err)
    }
    reader := bufio.NewReaderSize(logs, 64*1024)
    for next := true; next; {
        var l, line []byte

        for isPrefix := true; isPrefix; {
            l, isPrefix, err = reader.ReadLine()
            if err != nil {
                if err != io.EOF {
                    log.Fatal("Failed to retrieve logs", err)
                }
                next = false
            }
            line = append(line, l...)
        }
        if next || len(line) > 0 {
            var l logMessage
            err := json.Unmarshal(line, &l)
            if err != nil {
                //fmt.Printf("%+v", l)
                continue
            } else {
                if l.Change != nil {
                    if l.Level == "info" {
                        color.Green(string(l.Message))
                    } else {
                        color.Red(string(l.Message))
                    }
                }
                if l.Type == "change_summary" {
                    color.Magenta(string(l.Message))
                }
            }
        }
    }
    fmt.Println("")
}

type TerraformConfig struct {
    Credentials struct {
        Host struct {
            Token string `json:"token"`
        } `json:"app.terraform.io"`
    } `json:"credentials"`
}

func (tc *TerraformConfig) GetToken() (string, error) {
    token := os.Getenv("TFE_TOKEN")
    if token != "" {
        return token, nil
    }
    userHomeDir, err := os.UserHomeDir()
    if err != nil {
        log.Fatal(err)
    }

    content, err := os.ReadFile(filepath.Join(userHomeDir, "/.terraform.d/credentials.tfrc.json"))
    if err != nil {
        log.Fatal("🔐", "warning", "to authenticate set TFE_TOKEN or run terraform login")
    }

    err = json.Unmarshal(content, &tc)
    if err != nil {
        log.Fatal("Error when opening file: ", err)
    }
    return tc.Credentials.Host.Token, err
}

func Client() (*tfe.Client, error) {
    t := TerraformConfig{}
    token, err := t.GetToken()
    if err != nil {
        log.Fatal("Had a problem getting your API token", err)
    }
    config := &tfe.Config{Token: token}
    client, err := tfe.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }
    return client, err
}

// DownloadFile will download an url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
// Because DownloadURL is just a 302 to https://archivist.hashicorp.com/ which requires a custom header for auth
func downloadStateFile(filePath string, url string) error {
    t := TerraformConfig{}
    token, err := t.GetToken()
    if err != nil {
        return fmt.Errorf("had a problem getting your API token: %w", err)
    }

    client := retryablehttp.NewClient()
    client.RetryMax = 5
    client.Logger = nil
    // Set custom headers
    req, err := retryablehttp.NewRequest("GET", url, nil)
    if err != nil {
        fmt.Println("Error creating request:", err)
        return err
    }
    req.Header.Set("Authorization", "Bearer "+token)

    // Create the file
    out, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer out.Close()

    // Send the request
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    // Check server response
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("bad status: %s", resp.Status)
    }
    // Writer the body to file
    _, err = io.Copy(out, resp.Body)
    if err != nil {
        return err
    }
    return nil
}

// downloadStruct will download a struct to a file
// must use jsonapi rather than json.Marshal because there are jsonapi.NullableAttr[] in the workspace struct
func downloadStruct(filePath string, data interface{}) error {
    tempFile, err := os.CreateTemp("", "tfc")
    if err != nil {
        return fmt.Errorf("could not create temp file: %w", err)
    }
    defer os.Remove(tempFile.Name())
    defer tempFile.Close()

    if err := jsonapi.MarshalPayload(tempFile, data); err != nil {
        return fmt.Errorf("could not marshal data: %w", err)
    }
    // Need to do this in two steps to write indented json
    // Wish I could do it in one step
    tempFile.Seek(0, 0)
    var jsonData interface{}
    if err := json.NewDecoder(tempFile).Decode(&jsonData); err != nil {
        return fmt.Errorf("could not decode temp file: %w", err)
    }

    finalFile, err := os.Create(filePath)
    if err != nil {
        return fmt.Errorf("could not create file: %w", err)
    }
    defer finalFile.Close()

    encoder := json.NewEncoder(finalFile)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(jsonData); err != nil {
        return fmt.Errorf("could not encode to file: %w", err)
    }

    return nil
}
