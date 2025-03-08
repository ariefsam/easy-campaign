package tiktokapi

type tiktokApi struct{}

func New() *tiktokApi {
	return &tiktokApi{}
}

/*
{"userInfo":{"user":{"id":"6881290705605477381","shortId":"","uniqueId":"taylorswift","nickname":"Taylor Swift","avatarLarger":"https://p16-sign-va.tiktokcdn.com/tos-maliva-avt-0068/f8c7cad75f3a12205d31d6662d2555d5~tplv-tiktokx-cropcenter:1080:1080.jpeg?dr=14579&nonce=27329&refresh_token=2b7da31f64cab456e108d1f5565afee2&x-expires=1741215600&x-signature=ZOdsoNHJzg%2F7AcpfRzAtfbtInyY%3D&idc=my&ps=13740610&shcp=81f88b70&shp=a5d48078&t=4d5b0474","avatarMedium":"https://p16-sign-va.tiktokcdn.com/tos-maliva-avt-0068/f8c7cad75f3a12205d31d6662d2555d5~tplv-tiktokx-cropcenter:720:720.jpeg?dr=14579&nonce=79864&refresh_token=afebb214a1d991c019768c732e455e99&x-expires=1741215600&x-signature=UFfW2n4nnstK1q6k6GW1RQJnP30%3D&idc=my&ps=13740610&shcp=81f88b70&shp=a5d48078&t=4d5b0474","avatarThumb":"https://p16-sign-va.tiktokcdn.com/tos-maliva-avt-0068/f8c7cad75f3a12205d31d6662d2555d5~tplv-tiktokx-cropcenter:100:100.jpeg?dr=14579&nonce=11843&refresh_token=dedff0627a53aafab4f0f5c6a4c3fd94&x-expires=1741215600&x-signature=U5xt0vqdZI2kUUowrBjgE%2B%2FqcW4%3D&idc=my&ps=13740610&shcp=81f88b70&shp=a5d48078&t=4d5b0474","signature":"This is pretty much just a cat account","createTime":1628514847,"verified":true,"secUid":"MS4wLjABAAAAqB08cUbXaDWqbD6MCga2RbGTuhfO2EsHayBYx08NDrN7IE3jQuRDNNN6YwyfH6_6","ftc":false,"relation":0,"openFavorite":true,"bioLink":{"link":"taylorswift.com","risk":0},"commentSetting":0,"commerceUserInfo":{"commerceUser":false},"duetSetting":0,"stitchSetting":0,"privateAccount":false,"secret":false,"isADVirtual":false,"roomId":"","uniqueIdModifyTime":0,"ttSeller":false,"region":"US","downloadSetting":0,"profileTab":{"showMusicTab":true,"showQuestionTab":false,"showPlayListTab":false},"followingVisibility":1,"recommendReason":"","nowInvitationCardUrl":"","nickNameModifyTime":1629723423,"isEmbedBanned":false,"canExpPlaylist":true,"profileEmbedPermission":1,"language":"en","eventList":[],"suggestAccountBind":false,"isOrganization":0},"stats":{"followerCount":32700000,"followingCount":0,"heart":248100000,"heartCount":248100000,"videoCount":71,"diggCount":2204,"friendCount":0},"itemList":[]},"shareMeta":{"title":"Taylor Swift on TikTok","desc":"@taylorswift 32.7m Followers, 0 Following, 248.1m Likes - Watch awesome short videos created by Taylor Swift"},"statusCode":0,"statusMsg":"","needFix":true}
*/
type User struct {
	UserInfo struct {
		User struct {
			ID           string `json:"id"`
			ShortID      string `json:"shortId"`
			UniqueID     string `json:"uniqueId"`
			Nickname     string `json:"nickname"`
			AvatarLarger string `json:"avatarLarger"`
			AvatarMedium string `json:"avatarMedium"`
			AvatarThumb  string `json:"avatarThumb"`
			Signature    string `json:"signature"`
			CreateTime   int64  `json:"createTime"`
			Verified     bool   `json:"verified"`
			SecUID       string `json:"secUid"`
			Ftc          bool   `json:"ftc"`
			Relation     int    `json:"relation"`
			OpenFavorite bool   `json:"openFavorite"`
			BioLink      struct {
				Link string `json:"link"`
				Risk int    `json:"risk"`
			} `json:"bioLink"`
			CommentSetting   int `json:"commentSetting"`
			CommerceUserInfo struct {
				CommerceUser bool `json:"commerceUser"`
			} `json:"commerceUserInfo"`
			DuetSetting        int    `json:"duetSetting"`
			StitchSetting      int    `json:"stitchSetting"`
			PrivateAccount     bool   `json:"privateAccount"`
			Secret             bool   `json:"secret"`
			IsADVirtual        bool   `json:"isADVirtual"`
			RoomID             string `json:"roomId"`
			UniqueIDModifyTime int64  `json:"uniqueIdModifyTime"`
			TtSeller           bool   `json:"ttSeller"`
			Region             string `json:"region"`
			DownloadSetting    int    `json:"downloadSetting"`
			ProfileTab         struct {
				ShowMusicTab    bool `json:"showMusicTab"`
				ShowQuestionTab bool `json:"showQuestionTab"`
				ShowPlayListTab bool `json:"showPlayListTab"`
			} `json:"profileTab"`
			FollowingVisibility    int        `json:"followingVisibility"`
			RecommendReason        string     `json:"recommendReason"`
			NowInvitationCardUrl   string     `json:"nowInvitationCardUrl"`
			NickNameModifyTime     int64      `json:"nickNameModifyTime"`
			IsEmbedBanned          bool       `json:"isEmbedBanned"`
			CanExpPlaylist         bool       `json:"canExpPlaylist"`
			ProfileEmbedPermission int        `json:"profileEmbedPermission"`
			Language               string     `json:"language"`
			EventList              []struct{} `json:"eventList"`
			SuggestAccountBind     bool       `json:"suggestAccountBind"`
			IsOrganization         int        `json:"isOrganization"`
		} `json:"user"`
		Stats struct {
			FollowerCount  int `json:"followerCount"`
			FollowingCount int `json:"followingCount"`
			Heart          int `json:"heart"`
			HeartCount     int `json:"heartCount"`
			VideoCount     int `json:"videoCount"`
			DiggCount      int `json:"diggCount"`
			FriendCount    int `json:"friendCount"`
		} `json:"stats"`
		ItemList []struct{} `json:"itemList"`
	} `json:"userInfo"`
	ShareMeta struct {
		Title string `json:"title"`
		Desc  string `json:"desc"`
	} `json:"shareMeta"`
	StatusCode int    `json:"statusCode"`
	StatusMsg  string `json:"statusMsg"`
	NeedFix    bool   `json:"needFix"`
}
