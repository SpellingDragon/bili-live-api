package dto

import (
	"fmt"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func TestUnmarshall(t *testing.T) {
	data := "{\"action\":\"投喂\",\"bag_gift\":null,\"batch_combo_id\":\"batch:gift:combo_id:76452112:2073012767:31036:1686404457.6671\",\"batch_combo_send\":{\"action\":\"投喂\",\"batch_combo_id\":\"batch:gift:combo_id:76452112:2073012767:31036:1686404457.6671\",\"batch_combo_num\":1,\"blind_gift\":null,\"gift_id\":31036,\"gift_name\":\"小花花\",\"gift_num\":1,\"send_master\":null,\"uid\":76452112,\"uname\":\"某角落公司财务统计员\"},\"beatId\":\"\",\"biz_source\":\"live\",\"blind_gift\":null,\"broadcast_id\":0,\"coin_type\":\"gold\",\"combo_resources_id\":1,\"combo_send\":{\"action\":\"投喂\",\"combo_id\":\"gift:combo_id:76452112:2073012767:31036:1686404457.6651\",\"combo_num\":1,\"gift_id\":31036,\"gift_name\":\"小花花\",\"gift_num\":1,\"send_master\":null,\"uid\":76452112,\"uname\":\"某角落公司财务统计员\"},\"combo_stay_time\":5,\"combo_total_coin\":100,\"crit_prob\":0,\"demarcation\":1,\"discount_price\":100,\"dmscore\":72,\"draw\":0,\"effect\":0,\"effect_block\":0,\"face\":\"https://i1.hdslb.com/bfs/face/4add3acfc930fcd07d06ea5e10a3a377314141c2.jpg\",\"face_effect_id\":0,\"face_effect_type\":0,\"float_sc_resource_id\":0,\"giftId\":31036,\"giftName\":\"小花花\",\"giftType\":0,\"gold\":0,\"guard_level\":0,\"is_first\":true,\"is_join_receiver\":false,\"is_naming\":false,\"is_special_batch\":0,\"magnification\":1,\"medal_info\":{\"anchor_roomid\":0,\"anchor_uname\":\"\",\"guard_level\":0,\"icon_id\":0,\"is_lighted\":0,\"medal_color\":1725515,\"medal_color_border\":12632256,\"medal_color_end\":12632256,\"medal_color_start\":12632256,\"medal_level\":23,\"medal_name\":\"同事萌\",\"special\":\"\",\"target_id\":2073012767},\"name_color\":\"\",\"num\":1,\"original_gift_name\":\"\",\"price\":100,\"rcost\":22251324,\"receive_user_info\":{\"uid\":2073012767,\"uname\":\"美月もも\"},\"remain\":0,\"rnd\":\"2192488585708463104\",\"send_master\":null,\"silver\":0,\"super\":0,\"super_batch_gift_num\":1,\"super_gift_num\":1,\"svga_block\":0,\"switch\":true,\"tag_image\":\"\",\"tid\":\"2192488585708463104\",\"timestamp\":1686404457,\"top_list\":null,\"total_coin\":100,\"uid\":76452112,\"uname\":\"某角落公司财务统计员\"}"
	g := &Gift{}
	jsoniter.UnmarshalFromString(data, g)
	fmt.Printf("%+v", g)
}
