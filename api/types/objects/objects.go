package objects

type Clan struct {
	Id    string  `json:"id"`
	Name  string  `json:"name"`
	Color string  `json:"color"`
	Room  float64 `json:"room"`
}

type Misc struct {
	Opened   float64 `json:"opened"`
	Messages float64 `json:"messages"`
}

type Settings struct {
	Friends  string `json:"friends"`
	Requests string `json:"requests"`
}

type User struct {
	Id         int            `json:"id"`
	Username   string         `json:"username"`
	Created    float64        `json:"created"`
	Modified   float64        `json:"modified"`
	Avatar     string         `json:"avatar"`
	Banner     string         `json:"banner"`
	Badges     []string       `json:"badges"`
	Blooks     map[string]int `json:"blooks"`
	Tokens     float64        `json:"tokens"`
	Perms      []string       `json:"perms"`
	Clan       Clan           `json:"clan"`
	Role       string         `json:"role"`
	Color      string         `json:"color"`
	Exp        float64        `json:"exp"`
	Inventory  []string       `json:"inventory"`
	Misc       Misc           `json:"misc"`
	Friends    []int          `json:"friends"`
	Blocks     []int          `json:"blocks"`
	Claimed    string         `json:"claimed"`
	Settings   Settings       `json:"settings"`
	Otp        bool           `json:"otp"`
	MoneySpent float64        `json:"moneySpent"`
}

type BasicUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Color    string `json:"color"`
	Avatar   string `json:"avatar"`
	Banner   string `json:"banner"`
}

type Message struct {
	Id      int     `json:"id"`
	User    User    `json:"user"`
	Content string  `json:"content"`
	Date    float64 `json:"date"`
	Deleted bool    `json:"deleted"`
	Edited  bool    `json:"edited"`
}

type ChatData struct {
	Error   bool    `json:"error"`
	Event   string  `json:"event"`
	Author  User    `json:"author"`
	Message Message `json:"message"`
}

type BazaarItem struct {
	Id     int     `json:"id"`
	Item   string  `json:"item"`
	Price  float64 `json:"price"`
	Seller string  `json:"seller"`
	Date   float64 `json:"date"`
}
