package main

import (
	"fmt"
        "strings"
        g "github.com/sergle/go-gerrit"
        "github.com/sergle/go-gerrit/change"
)

const ConfigFile = "gerrit.conf"


var gerrit *g.Gerrit
var theme *ColorTheme
var proj_alias map[string]string
// default limit of concurrent connections
var Concurrent_GETs int = 10

func main() {

    cfg, err := ReadConfig(ConfigFile)
    if err != nil {
        fmt.Printf("Error reading config - %s\n", err)
        return
    }

    gerrit = g.New(cfg.Gerrit.User, cfg.Gerrit.Password, cfg.Gerrit.Host, cfg.Gerrit.CI)

    if cfg.Gerrit.Connections > 0 {
        Concurrent_GETs = cfg.Gerrit.Connections
    }

    theme = NewColorTheme(&cfg.Color)

    // format:  alias:fullname
    proj_alias = make(map[string]string)
    for _, alias := range cfg.Project.Alias {
        parts := strings.SplitN(alias, ":", 2)
        proj_alias[ parts[1] ] = parts[0]
    }

    dashboard("?q=owner:self+status:open")
    dashboard("?q=is:reviewer+status:open+-owner:self")
}

// get list of changes
func dashboard(query_string string) () {
    change_list, err := gerrit.FetchChangeList(query_string)
    if err != nil {
        return
    }

    num_changes := len(change_list)

    fmt.Printf("%sTotal %d changes%s\n", theme.Title, num_changes, theme.Reset)
    if num_changes == 0 {
        return
    }

    // get changes in parallel
    ch_out := make(chan *change.LongChange, num_changes)
    ch_in := make(chan string, num_changes)

    // limit concurrent get's 
    num_workers := Concurrent_GETs
    if num_workers > num_changes {
        num_workers = num_changes
    }

    // start 'workers'
    for i := 0; i < num_workers; i++ {
        go get_change(ch_in, ch_out)
    }

    for _, change := range change_list {
        ch_in <- change.Id
    }

    var processed = 0;
    ch_list := make([]*change.LongChange, 0)

    Loop:
    for {
        select {
            case change := <-ch_out:
                processed++
                ch_list = append(ch_list, change)
                if processed == num_changes {
                    break Loop
                }
        }
    }

    // and signal exit to worker goroutines
    for i := 0; i < num_workers; i++ {
        ch_in <- "EXIT"
    }

    // sort by Updated date
    gerrit.SortChanges(ch_list)
    print_change_list(ch_list)

    return
}

func get_change(ch_in <-chan string, ch_out chan<- *change.LongChange) () {

    for {
        select {
            case id := <-ch_in:

            if id == "EXIT" {
                return
            }

            detail, err := gerrit.GetChange(id)
            if err != nil {
                fmt.Printf("Failed to fetch change: %s\n", id)
            }

            ch_out <- detail
        }
    }

    return
}

func print_change_list(list []*change.LongChange) () {
    for _, ch := range list {

        ci_username, verified := gerrit.IsVerified(ch)

        var rating int8
        if verified < 0 {
            rating = -2
        } else {
            for _, p := range ch.Labels.CodeReview.All {
                if p.Value == -1 {
                    if rating > -1 {
                        rating = -1
                    }
                } else if p.Value == -2 {
                    rating = -2
                    break
                } else if p.Value == 1 && rating == 0 {
                    rating = 1
                } else if p.Value == 2 && rating >= 0 {
                    rating = 2
                }
            }
        }

        if ! ch.Mergeable {
            fmt.Printf("%s⨂%s ", theme.Bad, theme.Reset);
        } else if rating == 2 {
            fmt.Printf("%s✔%s ", theme.OK, theme.Reset);
        } else if rating == 1 {
            fmt.Printf("%s+%s ", theme.OK, theme.Reset);
        } else if rating == -1 {
            fmt.Printf("%s-%s ", theme.Bad, theme.Reset);
        } else if rating == -2 {
            fmt.Printf("%s✘%s ", theme.Bad, theme.Reset);
        } else {
            fmt.Printf("  ");
        }

        subj := ch.Subject
        if len(subj) > 50 {
            subj = subj[0:50]
        }

        proj, found := proj_alias[ch.Project]
        if ! found {
            proj = ch.Project
        }

        fmt.Printf("%-10s %-14s %-15s %-50s %s", proj, ch.Branch, ch.Owner.Username, subj, ch.Updated[0:16])

        ch_color := theme.Reset
        if verified == 1 {
            ch_color = theme.Verified
        } else if verified == -1 {
            ch_color = theme.Bad
        } else if verified == change.NotVerified {
            ch_color = theme.Absent
        }
        fmt.Printf(" %s%s%s", ch_color, ci_username, theme.Reset)

        for _, p := range ch.Labels.CodeReview.All {
            // exclude CI username (shown for verified)
            if gerrit.IsCI(p.Username) {
                continue
            }
            rv_color := theme.Reset

            if p.Value == -1 || p.Value == -2 {
                rv_color = theme.Bad
            } else if p.Value == 1 {
                rv_color = theme.Plus
            } else if p.Value == 2 {
                rv_color = theme.OK
            }

            fmt.Printf(" %s%s%s", rv_color, p.Username, theme.Reset)
        }

        fmt.Printf("\n")
    }
}

