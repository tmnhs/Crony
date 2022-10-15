package notify

var feiShuTemplateCard = `{
  "msg_type": "interactive",
  "card": {
    "config": {
      "wide_screen_mode": true
    },
    "header": {
      "title": {
        "tag": "plain_text",
        "content": "subjectSlot - Crony定时任务平台报警" 
      },
      "template": "red" 
    },
    "elements": [
				{
			  "fields": [
				{
				  "is_short": true,
				  "text": {
					"content": "**🕐 时间：**\ntimeSlot",
					"tag": "lark_md"
				  }
				},
				{
				  "is_short": true,
				  "text": {
					"content": "**📋报警主机：**\nipSlot",
					"tag": "lark_md"
				  }
				},
				{
				  "is_short": true,
				  "text": {
					"content": "**👤 值班：**\nuserSlot",
					"tag": "lark_md"
				  }
				},
				{
				  "is_short": false,
				  "text": {
					"content": "**报警信息:**\nmsgSlot",
					"tag": "lark_md"
				  }
				}
			  ],
			  "tag": "div"
			},
			{
			  "actions": [
				{
				  "tag": "button",
				  "text": {
					"content": "跟进处理",
					"tag": "plain_text"
				  },
				  "type": "primary",
				  "value": {
					"key1": "https://cloud.tencent.com/developer/article/1467743"
				  },
					"url":"https://cloud.tencent.com/developer/article/1467743"
				}
			  ],
			  "tag": "action"
			},
			{
			  "tag": "hr"
			}
    ]
  }
}`
