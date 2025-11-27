package main

import (
	"fmt"
	"github.com/lucky-lbc/jugglechat-server/utils"
)

func main() {
	var avatarUrls = []string{"https://pp-appicon.s3.amazonaws.com/avatar/default/default.png",
		"https://pp-appicon.s3.amazonaws.com/avatar/default/default.png",
		"https://pp-appicon.s3.amazonaws.com/avatar/default/default.png",
		"https://pp-appicon.s3.amazonaws.com/avatar/default/default.png",
		"https://pp-appicon.s3.amazonaws.com/avatar/default/default.png",
		"https://pp-appicon.s3.amazonaws.com/avatar/default/default.png",
		"https://pp-appicon.s3.amazonaws.com/avatar/default/default.png",
		"https://pp-appicon.s3.amazonaws.com/avatar/default/default.png",
	}
	err := utils.GenerateGroupAvatar(avatarUrls, "D:\\tmp\\im-server\\images\\group\\groups_cc.png")
	if err != nil {
		fmt.Println(err)
	}
}
