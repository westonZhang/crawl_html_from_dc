package main

import (
	"crawl_html_from_dc/services/get_response_html"
	"encoding/json"
	"fmt"
	"github.com/panwenbin/ghttpclient"
	"github.com/panwenbin/ghttpclient/header"
	"io/ioutil"
	"strings"
	"time"
)

type Send struct {
	Url    string `json:"url"`
	Header Header `json:"headers"`
}

type SendResponse struct {
	Msg string `json:"msg"`
}

type Receive struct {
	Url string `json:"url"`
}

type ReceiveResponse struct {
	Msg  string `json:"msg"`
	Err  string `json:"err"`
	Data struct {
		Status int    `json:"status"`
		Msg    string `json:"msg"`
		Rdata  string `json:"rdata"`
	} `json:"data"`
}

type Header struct {
	Cookie    string `json:"cookie"`
	UserAgent string `json:"user_agent"`
}

var receiveResults = make(map[string]*ReceiveResponse)

//var send_url = "http://127.0.0.1:9010/dc-send"
//var receive_url = "http://127.0.0.1:9010/dc-receive"

var send_url = "http://123.207.181.230:9010/dc-send"
var receive_url = "http://123.207.181.230:9010/dc-receive"

var headers = &Header{
	Cookie:    "BAIDUID=5028A9DDCD923E023F0FEEC0B22370EB:FG=1; BIDUPSID=5028A9DDCD923E023F0FEEC0B22370EB; PSTM=1568778486; BD_UPN=12314753; MSA_WH=1920_937; H_WISE_SIDS=135669_137150_137735_133103_136909_136651_136293_134725_113879_128065_136294_134982_136436_120195_137456_136659_137716_136366_132911_136455_135847_131247_137750_132378_131517_118881_118864_118849_118832_118788_136687_107319_132782_136799_136429_136091_133351_137222_136862_129649_136196_133847_132551_134047_131423_135232_136164_136753_110085_127969_131951_136612_137253_127416_136636_137097_137207_134349_132467_137619_137449_136987_100457; BDORZ=B490B5EBF6F3CD402E515D22BCDA1598; MCITY=-224%3A; delPer=0; BD_CK_SAM=1; BD_HOME=0; BDRCVFR[PowvNBqg-GC]=mbxnW11j9Dfmh7GuZR8mvqV; H_PS_PSSID=1461_21104_30073_29567_29221_26350_22159; PSINO=5; shifen[135882633018_59048]=1573543069; BCLID=11134997986530099843; BDSFRCVID=B7_OJeC627-0XRTw1RSOuQfJBejzquvTH6ao_9S7Ug_hlypMK8FoEG0PHx8g0Kub2p1TogKKL2OTHmuF_2uxOjjg8UtVJeC6EG0Ptf8g0M5; H_BDCLCKID_SF=tJAHoKLaJC83H43TqRrEKtFD-frQ5C62aKDsQpQ7BhcqEIL406Ah2xku5NLJQJotLCcZ5tjafR7kMxbSj4QohtAJ5Goi2t4OW6kJ2hoX3p5nhMJS257JDMP0-xQEXqQy523i2IovQpnVfqQ3DRoWXPIqbN7P-p5Z5mAqKl0MLPbtbb0xXj_0D6J3eaLHJ58s56bL3RTsH4jaKROvhDTjh6PYjnn9BtQmJJufsCJ9LfbbhfobXnoGbxIYbf6EbRQqQg-q3R77fx8bSJ33M-vBKMuUe-jy0x-jLgbOVn0MW-5Dh4tl3-nJyUPTD4nnBPrt3H8HL4nv2JcJbM5m3x6qLTKkQN3T-PKO5bRu_CcJ-J8XhDL4D5JP; H_PS_645EC=26d96KiTpOumgRlNnVQgkrrdhp69Dbb8Zrv8gqgnAMoQmglfLdlxs8xPebw; BDSVRTM=79; COOKIE_SESSION=3694_4_8_9_7_17_0_2_8_5_1_4_3461_0_0_0_1573540281_1573543069_1573549657%7C9%236730_3_1573543069%7C2; WWW_ST=1573551101383",
	UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.87 Safari/537.36",
}

var keywords3 = []string{
	"铁路运输代理服务",
	//"监控销售及安装", "建筑工程咨询", "销售不锈钢设备", "计算机网络技术", "新型建筑技术开发", "灌溉排水设施建筑", "电器自动化安装", "企业画册", "建筑机具租赁", "雕刻字", "电子监控设备的安装", "商务服务业", "安全帽", "民用刀貝及附件生产", "电子技术", "木托盘批发", "山水彩绘服务", "土地整理及土石方开挖", "商品混凝土研发", "机械设备的安装与维修", "桥梁建筑工程的承揽", "纺织配件", "矿石开采工程", "轻钢龙骨批发", "户外拓展活动策划", "通信设备安装工程", "电子设备的维修", "诱捕器具", "环评报告相关业务的代理服务", "从事城市轨道交通设备", "批发与零售物流设备", "管道工程的设计及施工", "工程竣工资料", "火灾自动报警系统", "预包装食品兼散装食品批发", "矿山机械设备销售", "有色金属技术研发", "水泥稳定层", "防水材料销售", "信息科技", "城市广场景观工程", "钢结构的制作及安装", "水果销售", "石灰石销售", "机械设备租赁和维修", "复印机", "石油工程技术研究", "医疗器械销售", "软件维护服务", "报验", "数码", "河流", "草料种植销售", "石油机械及配件", "广告牌设计", "信息化工程设计", "人才推荐", "不锈钢制品销售", "计算机信息科技技术服务", "水暖器材销售及安装", "装饰品销售", "数据存储处理", "商品网上销售", "公路工程的施工", "压力管道的安装", "处置废矿物油HW08", "污水处理技术开发", "农业科技的技术开发", "企业营销及管理咨询", "强弱电工程的设计", "防腐保温", "环保仪器设备", "生活垃圾", "防汛抗旱物资", "门窗设计", "装饰材料批发", "柴油机", "玉米秸", "窗帘制作", "网络维护", "市场信息咨询", "薯类", "消防金属制品", "生活污水拉运", "钾肥销售", "维修及技术咨询服务", "培育新品种及信息咨询服务", "洗井工程技术服务", "电气焊服务与维修", "人防工程工程", "热力设备清洁及安装", "防盗报警器材", "砌筑工程", "特种设备", "自营或代理进出口业务", "配电柜装置", "钢结构工程的加工", "劳动力租赁", "制作广告", "电子设备安装", "沙产业开发", "鞋帽零售", "室内外裝饰工程", "网络设备及计算机数码产品的安装", "农业机械设备研发", "土石方开挖回填工程", "2.4滴丁酯", "一二类医疗器械", "牛羊棚加工建设", "市政管理咨询服务", "网架工程", "氧焊", "桥梁架设", "智能机器人技术研发与销售", "房屋租赁及物业管理", "建筑材料的研发与技术服务", "电解铝", "工程项目管理及相关技术咨询", "乙酸溶液", "钢结构安装工程设计与施工", "水电暖设备销售及安装", "大气环境污染防治服务", "水利设施", "金属及非金属材料", "工程吊装设备租赁", "水土保持技术咨询服务", "超声检测管", "压力容器管道的安装", "建筑装饰装修及外墙保温工程承包", "喷涂设备及配件", "钛材制品", "机械设备批发兼零售", "道路修缮工程", "门禁考勤系统工程", "皮带机滚筒包胶", "从事新能源汽车科技领域内的技术开发", "铁线", "建筑机械及配件的销售", "计算机网络工程技术开发及维修", "第四代", "河湖水系工程", "不得从事危险化学品", "钻采设备研发", "编织", "计算机信息系统集成服务", "电气电子设备及元器件", "白渣灰", "果皮箱", "混凝土路面工程", "大型钢结构", "多媒体制作", "网络信息技术服务", "变频控制柜", "市政工程的技术咨询", "建筑工程机械设备进出口", "数据处理与存储服务", "电力工程设备安装", "仪器仪表及配件", "铁精矿粉", "船舶用品", "报刊图书", "窗", "机械备件", "艺术节", "机械设备及", "体外诊断试剂的批发及零售", "清洁保洁服务", "桥梁工程的施工", "设备存放服务", "水土流域综合治理", "桥梁加固工程", "互联网科技", "消防车辆", "闭杯点≤60°C", "花纹板", "全球导航卫星系统接收机", "不锈钢制品的销售", "建材运输服务", "美容仪器", "通信工程的设计与施工", "体育赛事活动策划服务", "非金属矿及制品经销", "市政园林环保工程", "塑钢及金属门窗设计", "网络文化产品展览", "文化艺术交流活动", "土地治理工程", "通信技术推广服务", "装修工程设计与施工", "种子批发及零售", "其他贸易经纪与代理", "家具零售", "农用设备", "自营和代理货物及技术进出口", "室内外装饰装修工程设备", "生物有机化肥", "环卫设备销售", "航空航天设备零部件", "钢板网加工及销售", "防漏", "房屋修缮工程的施工", "消防设备销售", "日用百货及办公用品", "勘察工程", "地坪工程", "防腐防渗工程", "工业自动化控制的设计", "室内外装修工程的设计", "其他厂房及建筑物工程建筑活动", "农业测土配方服务", "视频会议系统", "计算机软硬件开发与维护", "实验", "交通安全标识牌的安装", "室内装潢工程", "洗涤机械设备", "消防技术咨询服务", "及接插件生产", "房屋建设工程施工承包", "道路客运", "农机具维修", "城乡生活垃圾", "防虫网", "会议系统设备", "污水处理工程的设计施工", "建材的批发及零售", "管道铺设工程", "冲施肥", "劳保五金用品", "涉案财产", "普通道路货物运输服务", "大块石", "为农作物提供病虫害的防治服务", "普通机械设备销售", "标识标牌制作及销售", "纺织机械及配件", "从事智能科技", "软件系统开发", "堤坝工程", "小型工程", "掺混肥生产", "墙壁开关", "仿古园林建筑工程", "互联网服务", "施肥", "创业孵化咨询", "铁路机车配件", "珠宝玉器的销售", "装修及外保温工程承包", "国内劳务输出服务", "建筑设备及工程机械租赁", "新能源技术服务", "燃气设备销售", "销售汽车零配件", "园林古建筑工程施工", "水泥的生产", "及产品销售", "预制建筑物", "沙漠土地防治", "工程机械设备租赁及配件加工", "标志牌的设计与安装", "电线电缆和金属材料的销售", "低碳建筑的研发", "住宅装饰", "福美双乳剂", "成人用品", "提供农机作业服务", "技术交流和咨询服务", "液压管", "道路分包工程", "票务服务", "道路安全基础设施产业化技术研发", "二手机床销售服务", "机械设备及相关材料", "矿山物资销售", "装饰装修工程设计及施工", "输配电成套设备", "土石方开挖及回填工程", "二手车中介服务", "网络设备回收", "装饰装潢用品", "楼宇清洁服务", "光伏产品", "油田物资的技术服务及维修", "鞋", "卫生洁具销售", "塔吊租赁", "节水灌溉设备", "展会布置", "游艺器材及娱乐用品", "商品混泥土加工生产销售", "农药化肥", "彩钢房制作安装", "高低压配电器", "研究", "石材雕刻", "发展林下种植", "水电安装工程设备", "许可证有效期至2022年1月10日", "一般工业固体废物收集", "建材工程项目的评估", "家谱的设计与制作", "电子科技术开发", "电子元器件批发", "环评工程", "专项整理", "桥梁隧道施工设备的制作", "自动售货机", "接受金融机构委托从事金融信息技术外包", "电子设备的销售", "分切设备", "农业种植技术研发及推广服务", "园艺绿化服务", "舞台音响灯光设备", "混凝土工程的施工", "节水管理与技术咨询服务", "废旧机械设备折除", "推广服务及推广", "园林机械及", "土矿工程", "无人驾驶技术开发", "家政服务服务", "卫浴器材", "新能源", "建筑机械销售", "农林机械", "安全防范工程施工", "工程机械设备及配件的维修", "厨房及卫浴用具", "工艺品销售", "影音设备", "消防设备的销售", "节能设备改造", "商厦房屋", "服装批发服装", "压力容器焊接与安装工程", "厨房用具及日用杂品", "矿用机械设备及配件租赁", "水文地质勘察服务", "内胎", "园林设计", "舞台设备租赁", "环保工程技术的研发", "幕墙工程施工", "节能技术开发", "建筑技术咨询", "热处理设备", "水泥及制品", "种子销售", "土木", "电子工业产品", "水电暖设备及材料", "自动化设备销售", "商品信息咨询服务", "电话催缴", "网络约车服务", "道路园林工程承揽", "砂石的回收", "建筑安装工程承揽", "劳务承揽及施工", "建筑物清洗", "农业技术开发技术转让", "皮革皮具", "维修工程机械设备", "混凝土设备租赁", "水电安装工", "电器配件销售", "园林绿化及病虫害防治技术服务", "水库清淤", "锅炉安装销售", "热轧卷板", "桥梁涵洞工程", "地质勘查", "特种工程的设计", "食用农产品销售", "钢材零售", "活动房", "大型机械设备", "采煤机挖机配件", "农副产品初级加工", "应取得相关", "土体整理服务", "钻釆设备及配件", "非自有房屋租赁服务", "干混砂浆", "网络系统工程", "城市自来水和热力管道焊接", "绿化植物的种植及销售", "代收代缴电话费", "苗木花卉的种植", "工业化自动控制设备的销售", "自动化设备设计", "乙酸乙烯酯", "休闲农庄与观光旅游及餐饮服务", "货车", "劳务防护用品", "铁丝", "木制品包装", "电控", "UPS不间断电源","房屋代理销售", "环保填弃工程", "新能源技术研发", "阀", "日用陶瓷", "柔性铸铁排水管", "房屋建筑工程设计服务", "草原灭鼠", "防雷工程的施工", "设计及咨询服务", "道路客运服务", "照明及亮化器材", "消防工程服务", "电子科技的技术服务咨询", "焊接修理服务", "矿渣", "房屋建筑业", "石方工程施工服务", "面制品", "五金机电销售", "毒杂草防治及鼠虫害治理", "企业自产的商品", "太阳能设备安装及维护", "生态环境", "机电设备的技术服务", "电动汽车技术研发", "餐饮技术开发", "机电设备及零件", "养护服务", "工矿建筑工程", "生物科技领域内的技术开发", "文体及办公用品", "灯箱", "土石方开挖及回填", "橱柜", "水泥免烧砖", "沼泽湿地保护", "电子电气产品的生产", "厨电卫浴", "油水气分离设备", "注塑机维修", "石油钻采专用设备制造及维修", "舞台造型策划", "仓储货物运输", "独轮车的销售", "利用与购销", "道路施工的工程", "外墙清洗服务", "大豆油", "机电产品的批发零售及进出口", "邮票", "液压件", "城市生产生活垃圾", "污水处理设备的销售", "建筑生产", "机电安装工程施工", "门窗工程设计及安装", "机械租赁及运输", "摩托车及配件的销售", "农产品生产技术咨询服务", "须", "茶艺服务", "砂石料收购销售", "工程无纺布", "水电安装工程专业承包", "垃圾清运服务", "起重工程机械设备设施租赁", "折弯", "煤炭的销售", "起重机械及配件", "配套机电安装工程", "防水材料导入设备", "型材批发", "地基与基础工程承揽", "电子数码产品及电子元器件和家用电器的销售与维修", "建筑工程装潢材料的设计和施工", "节能节水设备", "室内外装潢及设计", "铁路运输设备租赁服务", "岩土工程勘察服务", "市政工程的施工及劳务分包", "各类造林绿化及城镇绿化苗木种植", "智能安防设备", "刻字刀", "建筑劳务人员派遣服务", "空气净化器租赁服务", "园林绿化及景观工程", "繁肓", "苗木花卉销售", "生活污水处理服务", "窗制作", "市政管道工程", "视频监控存储设备", "儿童玩具的批发及零售", "水电工程前期报告的编制和咨询", "汽车维修与保养", "家电设备", "非自有房屋租赁", "水果种植及销售", "网架及配件", "橡胶及其制品", "包装装潢", "砖石砌筑", "电气机械设备", "城市停车场服务", "公路路基工程施工", "环卫设施", "工程勘察活动", "建筑土建施工", "装饰材料的", "网围栏的销售及安装", "石油和天然气开采专业及辅助性活动", "金刚石工具", "人才中介服务", "水平衡测试服务", "草", "工业自动化系统装置", "苗木收购", "设计和维修", "绿化工程设计及施工", "餐饮企业管理服务", "新材料技术开发", "房屋建筑工程及设计服务", "LED电子显示屏设计与制作", "建筑装饰设备的销售和租赁", "灌溉服务活动", "矿山工程技术咨询服务", "水文地质", "工业技术支持及咨询服务", "酒店开发及餐饮服务", "建筑装饰工程的招标代理", "机械设备及零配件的销售", "机械设备及场地的租赁", "户外运动装备", "粮食仓储服务", "大型工程机械设备发动机", "人工智能科技领域内技术的研发", "速冻食品", "的设计制造与安装咨询服务", "工程机械车辆租赁", "彩钢房搭设", "体育拓展活动策划", "仿古建筑工程施工", "变性淀粉", "屋面防水作业", "其他仓储业", "电脑图文制作", "山产品", "轻小型起重设备", "监测服务", "服装销售", "化学制剂", "求援", "林业有害生物防治", "企业管理信息咨询服务", "电梯安装搭设拆除", "马铃薯", "光缆电缆", "铝合金玻璃", "固井工程及相关技术服务", "工程机械的批发与零售", "野营房的制造及销售", "计算机软硬件及辅助设备耗材", "汽车及汽车零部件", "通讯设备及相关产品", "电热元件", "新技术服务", "仪器仪表租赁", "公路交安工程", "计算机软硬件及辅助设备销售", "五金机械", "建筑机械设备及材料", "污泥处置", "供热工程", "楼宇清洗", "建筑防水材料的销售", "农业休闲观光活动", "通信终端设备", "森林防火标识牌", "输配电及成套电控装置", "市政道路工程服务", "豆类及薯类批发", "采矿设备", "路牌路标及广告牌", "牛羊肉及鲜活禽蛋", "代理及广告的推广", "智能设施设备", "农业旅游开发", "电管道工程的安装服务", "帽", "防静电设备", "电力安装工程劳务分包", "消防设施工程的设计", "电机及控制系统", "建筑装饰装修工程施工与设计", "五金建材零售", "以本社成员为主要服务对象", "军用车零部件", "发布各类", "机电设备及配件安装与销售", "建筑装饰工程服务", "配电箱成套组装", "音响设备研发", "金属器材", "粉煤灰及工业废渣", "空气治理", "护理服务", "框架眼镜", "门式起重机", "钾肥包装", "智能建筑设计", "汽车钣金及装饰服务", "道路工程养护及绿化", "旅游咨询服务", "互联网搜索服务", "技术销售", "机动车新车上牌服务", "食品的批发及零售", "有色金属加工", "数控设备", "水电安装工程的设计与施工", "非金属及金属矿", "自动化仪表", "起吊", "洁具卫浴的批发", "市场营销策划及推广", "节能节水设备上门安装服务", "豆及薯类销售", "投资", "断桥铝门窗制作销售", "模具的加工", "防爆设备及配件的销售", "湿地养护", "工业污水", "道路工程施工及维修维护", "落户手续", "化妆品研发", "小修", "制冷设备的生产", "电力设施工程承接", "旅游地产开发", "公路工程材料的销售", "输配电设备", "农田基地建设服务", "室内外墙面粉刷", "铝型材销售", "农业科学技术的研发与推广", "电力科技的技术开发", "道路的维护", "交通安全设施工程", "物联网服务", "改造及销售", "景观雕塑工程", "硫酸钾肥", "吊装设备的租赁", "门窗及型材喷涂加工", "园林养护用品", "大型机械设备的拆解", "五金设备", "太阳能设备安装", "电暖及新能源材料代购代销", "除废电瓶拆解", "绿化栽植", "石材雕刻刀", "冰雪雕刻", "汽车租赁#", "玉制品", "厨房用具及日用杂品零售", "国内各类广告的设计", "电脑喷绘", "环保设施及产品销售", "饲料及饲料添加剂", "电站设备及配件", "编码器", "消防设备工程施工", "气缸床", "茶具销售", "工程劳务的承包", "环境信息咨询", "投影仪", "电子产品的租赁", "信息系统集成", "敌百虫", "家居装饰材料", "防水防腐保温工程专业承包", "施工及安装", "面制品及食用油的批发", "汽车清洗", "防水胶布", "管网工程的施工", "工程机械租赁及材料销售", "农机", "装饰材料的销售及安装", "金属门窗的加工", "机械设备的安装与租赁", "室内装饰材料销售", "机动车辆", "货物信息配载", "推广和销售", "电器设备运行维护", "水利水电工程的施工", "电信业务", "电力设备维修工程", "冷热水的运输及配送", "分离设备零配件",
	//"计算机科学技术研究服务", "含婴幼儿乳粉", "安装报警工程", "淡水鱼养殖", "安防设备销售及安装", "二灰石", "路基材料", "教育技术推广服务", "道路养护及维修", "随车工具", "农业机械配件销售", "网围栏配套产品加工", "钢材料的销售", "安装及售后服务", "环保产品批发零售", "新能源工程承揽", "标识牌", "书刊", "油井配套设备及相关材料", "智能产品销售", "农业技术的研发与技术咨询服务", "异型件的加工制造", "社会调研", "须经审批的事项", "园林绿化工程施工承包", "苗圃", "大型工程机械设备租赁", "针纺用品", "建筑用附件及机械租赁", "矿石加工", "展厅的布置设计", "化工能源设备", "建筑机械设备销售", "其他建筑安装业", "农副产品初加工及销售", "市政公共设施安装", "标准化紧固件", "其他日用品零售", "机械设备安装销售", "金属门窗和路灯的安装及销售", "铁路运输代理服务", "非标准及专用实验设备", "土石方工程承包", "马铃薯种植及销售", "水暖电料材料的销售", "文化艺术活动交流策划", "混凝土破碎", "销售环保材料", "贸易中介代理", "燃烧设备", "智能家居控制系统", "不锈钢型材", "监控系统安装及销售", "选矿药剂", "室内外装饰装潢工程设计施工", "娱乐", "数控机床设备", "机械设备及零件", "人工造林服务", "的零售", "工程机械服务", "乡村道路建设与维修工程", "钟表销售", "机械设备的销售及维修", "计量电度表", "人防工程专业承包", "母线槽", "报警设备", "制冷设备安装工程的施工", "技术指导及服务", "管道及配件安装", "绿色植物培植", "环境控制设备", "楼宇设备自控系统工程", "渣土及普通货物运输", "地基与基础的施工", "金属家具制作", "广告图文设计", "网围栏加工及销售", "安全防护用品", "工程机械出租", "土特产品加工", "以自有资金对医药行业进行投资", "建筑材料的生产", "维修与调试", "民用航空材料", "不得开设储煤场", "外墙涂料销售", "机电产品的制造", "放电记录器", "家用电品", "道路养护服务", "动力电池", "环保建材制造", "设备的维护", "卷板机", "桥梁附属结构制作", "电气设备安装服务", "塑钢门窗销售", "仪器和器材销售", "利用互联网销售农副产品", "环境工程设计服务", "装卸吊运", "建筑工程管理", "计算机现场维修", "美容美体", "家禽及销售", "石制建设工程作业", "电梯技术研发及技术咨询", "焊割设备", "资料服务", "结构件加工", "酒店餐饮服务", "监控及", "鲜大肉批发", "电力变压器", "园林绿化工程设备", "电动伸缩门安装及销售", "道路硬化工程承揽", "彩钢加工销售", "橡胶制品的生产", "瓜类", "交通管理有发光标志批发", "网络通讯工程", "土畜产品", "集装箱", "小麦草", "环境保护及生态治理工程", "装饰材料批发及零售", "油封", "天然砂开采", "易爆品", "杂粮种植", "路政施工工程", "电脑周边产品及辅助设备", "代收洗衣", "木料", "医疗器械的研发", "环保技术的技术研发", "林木种子生产", "热能设备及配件", "水利防汛物资", "项目咨询", "油漆化工", "苗木花卉的种植销售", "水泥预制工程", "网络的技术开发", "单板加工", "果品蔬菜", "锂电池材料", "石料的批发与销售", "的策划", "光纤或半导体激光打标机", "生产推广与销售", "机电设备上门维修", "建筑材料加工销售", "医疗科技技术开发", "以及动物和动物产品无害化处理场所", "浸塑护栏网", "监控安装工程的设计与施工", "水泥沙石运输及销售", "人工草坪与塑胶操场工程", "承办经批准的体育赛事活动", "建筑行业", "废旧机电机械设备回收", "售后维修", "电磁导热油锅炉", "汽车零部件的加工", "废电瓶", "市政公用设施销售", "从事电子科技领域内的技术开发", "多媒体设计服务", "其他日用品", "模具标准件", "林木种苗及草皮", "销售及维修保养", "防爆破", "仓库房屋工程建筑", "土方工程施工", "建筑铝合金模板的租赁", "电梯维修", "检查井砌块", "计算机软硬件及外围设备的技术开发", "农业科技技术咨询服务", "货物代理", "轮胎及橡胶制品的批发和零售", "土木建筑", "汽车螺丝", "甲醛溶液", "特种专业工程", "自动化控制系统及软件", "化妆护肤用品", "商务服务", "为餐饮企业提供送餐服务", "箱包配件", "混凝土浇灌", "废旧金属及有色金属的回收利用", "储罐机械清洗服务", "农业种植技术研发", "围栏销售与安装", "城市公路桥梁工程建筑活动", "货物中转", "防水工程的设计及施工", "保洁设备", "木材销售及加工", "智能设备的设计", "钢制暖气片制造", "绿化景观工程", "办公设备销售", "家俱", "土石石方工程", "场地出租", "保健酒", "公物", "砂纸砂布", "建筑防水防腐保温", "砌筑建设工程作业", "零售卷烟", "办公用品租赁", "金属系固体标准件管件", "机械设备销售维修", "公益活动策划", "砂石水泥销售", "洁具的销售", "室内环境检测与治理", "新能源材料", "专门为买卖双方提供贸易机会的中介代理", "废旧物资回收及销售", "橡胶材料销售", "绿化植物租赁", "工程机械及配件的维修", "音响设备租赁", "节能环保设备领域内的技术开发", "补光灯", "水暖管道安装工程", "科研所需的原辅材料", "机电一体化的技术及产品", "光卤石", "机械工具", "高空作业", "代办机动车产权过户手续", "房建及拆除’老旧楼房改建工程", "防洪除涝技术咨询服务", "新科技项目", "广告设计与制作", "液压元件", "搅拌器", "造林及更新工程", "饵料", "森林抚育", "草地治理", "压力传感器", "机械设备制造", "牧草等农作物种植", "农业技术开发和技术转让", "办公用品批发零售", "隧道维修", "声学装饰装修工程的设计施工", "建筑工程及技术咨询", "空气净化器", "综合布线系统工程", "互联网信息技术服务", "滤布滤袋", "水暖器材及设备", "家居用品设计", "建筑钢结构预制构件工程安装服务", "空调设备及配件", "建筑铝合金门窗和塑钢门窗加工", "物流仓储车", "旅游咨询", "机床刀具及量具", "环保科技领域内的技术研发", "景观土石方工程", "建筑工程质量检测", "商砼", "金属材料加工", "混泥土预制构件", "减速机", "矿山设备及工程机械", "通讯终端设备销售", "畜牧兽医设备", "广播电视节目制作", "金属制品及技术的出口业务和所需的机械设备", "拖拉机", "日杂品", "加碘食盐", "氧化镁", "灯具零售", "园林绿化技术咨询", "服装加工", "物联网智能综合平台的研发", "经销计算机软硬件及辅助设备", "收购初级农副产品", "金属栏杆设计安装",
}

func main() {
	//go send()

	i := 1
	for len(keywords3) > 0 {
		receive()
		unAccomplishKeywords := []string{}
		for keyword := range receiveResults {
			if receiveResults[keyword].Err != "" {
				if strings.Contains(receiveResults[keyword].Err, "下载中心的下载队列中暂时没有该任务") {
					if i >= 2 {
						reSend(keyword)
					}
				}
				unAccomplishKeywords = append(unAccomplishKeywords, keyword)
			} else {
				rDataMap, err := get_response_html.ResponseRDataMap(receiveResults[keyword].Data.Rdata)
				if err != nil {
					fmt.Println(err)
					continue
				}

				html, err := get_response_html.DecodeHtml(rDataMap)
				if err != nil {
					fmt.Println(err)
					continue
				}
				writeHtml(keyword, html)
			}
			time.Sleep(time.Millisecond * 200)
		}

		receiveResults = map[string]*ReceiveResponse{}
		keywords3 = unAccomplishKeywords
		fmt.Println("length unAccomplishKeywords:", len(unAccomplishKeywords), unAccomplishKeywords)
		i++
	}
}

func send() {
	for i, keyword := range keywords3 {
		sendData := &Send{
			Url:    fmt.Sprintf("http://www.baidu.com/s?wd=%s", keyword),
			Header: *headers,
		}

		jsonBytes, _ := json.Marshal(sendData)
		sendResponse := &SendResponse{}
		err := ghttpclient.PostJson(send_url, jsonBytes, nil).ContentType("application/json").ReadJsonClose(sendResponse)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(i, "/", len(keywords3), keyword, *sendResponse)
		time.Sleep(time.Millisecond * 200)
	}
}

func receive() {
	for _, keyword := range keywords3 {
		receiveData := &Receive{
			Url: fmt.Sprintf("http://www.baidu.com/s?wd=%s", keyword),
		}
		jsonBytes, _ := json.Marshal(receiveData)
		receiveResponse := &ReceiveResponse{}

		httpHeader := make(header.GHttpHeader)
		httpHeader.Set("Content-Type", "application/json;charset=UTF-8")
		err := ghttpclient.PostJson(receive_url, jsonBytes, httpHeader).ContentType("application/json").ReadJsonClose(receiveResponse)
		if err != nil {
			fmt.Println(keyword, err)
		}

		receiveResults[keyword] = receiveResponse
		fmt.Println("receiving result", keyword)
		time.Sleep(time.Millisecond * 200)
	}
}

func writeHtml(keyword string, html string) {
	path := fmt.Sprintf("./data/html/%s.html", keyword)
	data := []byte(html)
	if ioutil.WriteFile(path, data, 0644) == nil {
		fmt.Println("写入文件成功:", keyword)
	}
}

func reSend(keyword string) {
	sendData := &Send{
		Url:    fmt.Sprintf("http://www.baidu.com/s?wd=%s", keyword),
		Header: *headers,
	}

	jsonBytes, _ := json.Marshal(sendData)
	sendResponse := &SendResponse{}
	err := ghttpclient.PostJson(send_url, jsonBytes, nil).ContentType("application/json").ReadJsonClose(sendResponse)
	if err != nil {
		fmt.Println(err)
	}
}
