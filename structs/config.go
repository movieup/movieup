package structs

import (
)

// DBConfig represents database connection configuration information.
type GMailConfig struct { // toml内の名前を入れる
    MailAddress     string `toml:"MailAddress"`
    ClientID        string `toml:"ClientID"`
    ClientSecret    string `toml:"ClientSecret"`
    AccessToken     string `toml:"AccessToken"`
    RefreshToken    string `toml:"RefreshToken"`
}

type DBConfig struct { // toml内の名前を入れる
    User        string `toml:"user"`
    Password    string `toml:"password"`
    Protocol    string `toml:"protocol"`
    DBname      string `toml:"dbname"`
}

type EtcConfig struct { // toml内の名前を入れる
    SrcPath string `toml:"srcPath"`
}

// Config represents application configuration.
type Config struct { // toml内の名前を入れる
    Gmail   GMailConfig `toml:"gmail"`
    DB      DBConfig    `toml:"DB"`
    Etc     EtcConfig   `toml:"etc"`
}
