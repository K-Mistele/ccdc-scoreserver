package constants

import "time"

const UserAgent = "Blacklight-ScoreServer-1.0"
const BrowserFriendlyUserAgent = "Blacklight-ScoreServer 1.0, like: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:85.0) Gecko/20100101 Firefox/85.0"
var ServerTime, _ = time.LoadLocation("America/Chicago")
