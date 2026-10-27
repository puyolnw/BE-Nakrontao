-- Content drawn from the research proposal "การพัฒนาแอปพลิเคชันท่องเที่ยวสายมู เจ้าพ่อศรีนครเตา"
-- (Fundamental Fund 2569, มหาวิทยาลัยราชภัฏมหาสารคาม), section 5 "เจ้าพ่อศรีนครเตา" and section 6, which
-- cite: interview with community scholar จิรายุทธ แสงมณี (16 Apr 2567); the shrine history booklet handed out at the
-- shrine festival; Fine Arts Department survey of Thung Kula Ronghai 2529–2546; อมาวสี เถียรถาวร (2532);
-- บุญทัน กมล (2552); พิเนตร ดาวเรือง (2563).
-- Still to be confirmed by the community reviewer (survey item A5): Chinese/English translations, exact
-- coordinates of each site (only the village coordinate is given), ceremony dates, shop contacts, ritual steps.

INSERT INTO site_settings(key, value) VALUES
('site', '{"name":{"th":"เจ้าพ่อศรีนครเตา","zh":"西那空涛神","en":"Chao Pho Si Nakhon Tao"},
  "tagline":{"th":"แพลตฟอร์มท่องเที่ยวเชิงศรัทธา","zh":"信仰旅游平台","en":"A faith tourism guide"},
  "contact":"มหาวิทยาลัยราชภัฏมหาสารคาม · 043-725437","survey_url":"","hero_url":"https://images.unsplash.com/photo-1758346971581-81d68e57e9b3?auto=format&fit=crop&w=2000&q=75"}'),
('privacy', '{"th":"แพลตฟอร์มนี้ไม่เก็บชื่อหรือข้อมูลติดต่อของผู้ใช้ คำถามที่ส่งถึงผู้ช่วย AI ถูกลบข้อมูลส่วนบุคคลออกก่อนบันทึก และลบอัตโนมัติเมื่อครบ 180 วัน","zh":"本平台不收集用户姓名或联系方式。发送给AI助手的问题在保存前会去除个人信息，并在180天后自动删除。","en":"We do not collect names or contact details. Questions sent to the assistant are stripped of personal data before storage and deleted after 180 days."}'),
('consent', '{"th":"แบบสอบถามนี้เป็นส่วนหนึ่งของโครงการวิจัย ข้อมูลใช้เพื่อการวิเคราะห์ภาพรวมเท่านั้น ไม่ระบุตัวตน ท่านสามารถยุติการตอบได้ทุกเมื่อ","zh":"本问卷是研究项目的一部分，数据仅用于匿名的整体分析，您可随时停止作答。","en":"This questionnaire is part of a research project. Answers are anonymous and analysed in aggregate only. You may stop at any time."}');

-- ---------------------------------------------------------------- knowledge
INSERT INTO knowledge(doc_code, category, kind, title, body, source, informant, reviewer, status, sort) VALUES
('K-001','community','history',
 '{"th":"ชุมชนบ้านเมืองเตา","zh":"孟涛村社区","en":"The village of Ban Mueang Tao"}',
 '{"th":"บ้านเมืองเตาตั้งอยู่ในตำบลเมืองเตา อำเภอพยัคฆภูมิพิสัย จังหวัดมหาสารคาม ประกอบด้วยหมู่ที่ 1 หมู่ที่ 11 และหมู่ที่ 12 จากหมู่บ้านทั้งหมด 26 หมู่บ้านในตำบล อยู่ห่างจากตัวอำเภอประมาณ 20 กิโลเมตร มีพื้นที่กว่าหนึ่งหมื่นหกพันไร่ (ละติจูด 15.5116 ลองจิจูด 103.3866)\nพื้นที่อยู่สุดเขตด้านทิศตะวันออกของอำเภอ เป็นที่ราบชายทุ่งกุลาร้องไห้ ไม่มีภูเขา เป็นทุ่งนาสลับป่าโปร่งและป่าละเมาะ ดินปนทราย หน้าแล้งแล้งจัด หน้าน้ำน้ำท่วมทั่วทุ่ง\nทิศเหนือติดอำเภอปทุมรัตต์ จังหวัดร้อยเอ็ด ทิศตะวันออกติดอำเภอเกษตรวิสัย จังหวัดร้อยเอ็ด ทิศตะวันตกติดตำบลราษฎร์เจริญ และทิศใต้ติดอำเภอชุมพลบุรี จังหวัดสุรินทร์\nชาวบ้านส่วนใหญ่ทำนาเป็นหลัก รองลงมาคือรับจ้าง รอบหมู่บ้านมีสระน้ำโบราณหลายแห่ง และมีลำน้ำเตาไหลผ่านห่างจากชุมชนราว 2 กิโลเมตร",
   "zh":"孟涛村位于玛哈沙拉堪府帕亚卡普米披赛县孟涛乡，由第1、11、12村组组成（全乡共26个村），距县城约20公里，面积逾一万六千莱（北纬15.5116，东经103.3866）。\n村庄位于县境最东端，是古拉龙海平原边缘的平地，没有山丘，稻田与疏林交错，土壤为沙质，旱季极旱，雨季整片原野被水淹没。\n北接黎逸府巴吞叻县，东接黎逸府格塞威赛县，西接叻察仁乡，南接素林府春蓬武里县。\n村民以种稻为主，其次为打工。村子周围有多处古代水塘，涛河在距村约2公里处流过。",
   "en":"Ban Mueang Tao lies in Mueang Tao subdistrict, Phayakkhaphum Phisai district, Maha Sarakham. It covers villages no. 1, 11 and 12 of the 26 in the subdistrict, about 20 km from the district town, over 16,000 rai in area (15.5116 N, 103.3866 E).\nIt sits at the far eastern edge of the district on the flat margin of Thung Kula Ronghai, with no hills: rice fields among open woodland on sandy soil, parched in the dry season and flooded across the plain in the wet season.\nTo the north is Pathum Rat (Roi Et), to the east Kaset Wisai (Roi Et), to the west Rat Charoen subdistrict and to the south Chumphon Buri (Surin).\nMost villagers farm rice, others work for wages. Ancient ponds ring the village and the Tao stream flows about 2 km away."}',
 'ข้อเสนอโครงการ FF 2569 มรภ.มหาสารคาม หมวด 5','','รอผู้ตรวจสอบชุมชน (A5)','published',1),

('K-002','origin','history',
 '{"th":"ที่มาของชื่อ \"เมืองเตา\"","zh":"“孟涛”之名的由来","en":"Why the village is called Mueang Tao"}',
 '{"th":"ชื่อเมืองเตามาจากเตาถลุงเหล็ก ชุมชนตั้งอยู่ใกล้ลำน้ำเตา และมีหลักฐานการถลุงเหล็กสมัยก่อนประวัติศาสตร์ คือขี้เหล็กหรือตะกรันเหล็กกระจัดกระจายอยู่ทั่วบริเวณที่เรียกว่า \"หนองเตาเหล็ก\" ทางทิศตะวันตกเฉียงใต้ของหมู่บ้าน\nปัจจุบันหนองเตาเหล็กอยู่บริเวณหลังโรงเรียนเมืองเตาวิทยาคม ยังพบตะกรันเหล็กและศิลาแลงหลงเหลืออยู่มาก\nตามตำนาน เมื่อเซียงสีได้รับบรรดาศักดิ์เป็นหลวงศรีนครเตา บ้านกุดหวายจึงได้รับการยกฐานะเป็นบ้านเมืองเตา โดยคำว่า \"เตา\" หมายถึงเตาถลุงเหล็กสำหรับทำมีดและดาบ",
   "zh":"“涛”（Tao）意为炼铁炉。村庄靠近涛河，并有史前炼铁的证据：村西南方被称为“铁炉塘”（Nong Tao Lek）的地方散布着铁渣。\n如今铁炉塘位于孟涛威塔亚空中学后方，仍可见大量铁渣与红土石。\n据传说，当谢西（Xiang Si）受封为“銮·西那空涛”后，古外村被提升为孟涛村，“涛”即用来打造刀剑的炼铁炉。",
   "en":"Tao means an iron-smelting furnace. The village lies near the Tao stream, and prehistoric iron working left slag scattered across a place called Nong Tao Lek, the iron-furnace pond, to the south-west of the village.\nToday Nong Tao Lek is behind Mueang Tao Witthayakhom School, where slag and laterite are still plentiful.\nIn the legend, when Xiang Si was made Luang Si Nakhon Tao, the village of Kut Wai was raised to Ban Mueang Tao, tao being the furnace in which knives and swords were forged."}',
 'ข้อเสนอโครงการ FF 2569 หมวด 5; หนังสือประวัติเจ้าพ่อศรีนครเตา (แจกในวันฉลองศาล)','','รอผู้ตรวจสอบชุมชน (A5)','published',2),

('K-003','archaeology','history',
 '{"th":"แหล่งโบราณคดีในตำบลเมืองเตา","zh":"孟涛乡的考古遗址","en":"Archaeological sites in Mueang Tao subdistrict"}',
 '{"th":"การสำรวจแหล่งโบราณคดีพื้นที่ทุ่งกุลาร้องไห้ของกรมศิลปากรช่วง พ.ศ. 2529–2546 พบชุมชนโบราณกระจายอยู่ในตำบลเมืองเตา ได้แก่\n1. แหล่งโบราณคดีบ้านเมืองเตา ชุมชนโบราณมีคูน้ำคันดิน อายุก่อนประวัติศาสตร์ถึงสมัยประวัติศาสตร์วัฒนธรรมเขมร\n2. โนนหนองคอน (บ้านโคกสาย) 3. โนนน้อย (บ้านโนนยูง) 4. วัดป่าโนนยูง (บ้านโนนยูง) อายุก่อนประวัติศาสตร์ตอนปลายถึงวัฒนธรรมเขมร\n5. โนนน้ำอ้อมใหญ่ (บ้านโนนจาน) 6. บ้านน้ำอ้อม 7. โนนป่าตาล (บ้านน้ำอ้อม) ชุมชนโบราณมีคูน้ำคันดิน\n8. บ้านกางกี่ สมัยประวัติศาสตร์วัฒนธรรมเขมร 9. โนนโพนโล้น (บ้านกางกี่) มีคูน้ำคันดิน\n10. บ้านโนนบ่อ 11. โนนหว้า (บ้านโนนปอ) 12. โนนกกข่อย (บ้านยานาง) 13. โนนหนองผักสั่ง และ 14. บ้านผักสั่ง (บ้านผักสั่ง)\nแหล่งส่วนใหญ่มีอายุก่อนประวัติศาสตร์ตอนปลาย สอดคล้องกับวัฒนธรรมทุ่งกุลาร้องไห้ยุคโลหะตอนปลาย",
   "zh":"泰国艺术厅于佛历2529–2546年（1986–2003年）对古拉龙海平原进行考古调查，在孟涛乡发现多处古代聚落，包括：\n1. 孟涛村遗址：有护城河与土墙的古代聚落，年代从史前延续到高棉文化时期。\n2. 农农空（谷赛村）3. 农诺（农允村）4. 巴农允寺（农允村），史前晚期至高棉文化时期。\n5. 农南翁艾（农占村）6. 南翁村 7. 农巴丹（南翁村），均为有护城河土墙的古代聚落。\n8. 康基村（高棉文化时期）9. 农蓬伦（康基村），有护城河土墙。\n10. 农波村 11. 农华（农波村）12. 农谷开（雅囊村）13. 农农帕桑 及 14. 帕桑村。\n大多数遗址属史前晚期，与古拉龙海铁器时代晚期文化一致。",
   "en":"The Fine Arts Department surveyed Thung Kula Ronghai between 1986 and 2003 (BE 2529–2546) and recorded ancient settlements across Mueang Tao subdistrict:\n1. Ban Mueang Tao, a moated settlement with earthen banks, from prehistory into the Khmer period.\n2. Non Nong Khon (Ban Khok Sai), 3. Non Noi (Ban Non Yung), 4. Wat Pa Non Yung (Ban Non Yung): late prehistoric to Khmer.\n5. Non Nam Om Yai (Ban Non Chan), 6. Ban Nam Om, 7. Non Pa Tan (Ban Nam Om): moated settlements.\n8. Ban Kang Ki (Khmer period), 9. Non Phon Lon (Ban Kang Ki), moated.\n10. Ban Non Bo, 11. Non Wa (Ban Non Po), 12. Non Kok Khoi (Ban Ya Nang), 13. Non Nong Phak Sang and 14. Ban Phak Sang.\nMost date to late prehistory, matching the late Metal Age culture of Thung Kula Ronghai."}',
 'กรมศิลปากร สำรวจแหล่งโบราณคดีทุ่งกุลาร้องไห้ พ.ศ. 2529–2546 (อ้างใน ข้อเสนอโครงการ FF 2569)','','รอผู้ตรวจสอบชุมชน (A5)','published',3),

('K-004','archaeology','history',
 '{"th":"โบราณวัตถุและร่องรอยเมืองโบราณ","zh":"古代文物与古城遗迹","en":"Artefacts and traces of the ancient town"}',
 '{"th":"บ้านเมืองเตาตั้งอยู่บนที่ดอน มีคูเมืองขนาดใหญ่ล้อมรอบชั้นเดียว และมีสระน้ำโบราณรอบชุมชน ซึ่งเชื่อว่าผู้คนใช้ตั้งถิ่นฐานมาตั้งแต่ยุคโลหะตอนปลายในวัฒนธรรมทุ่งกุลาร้องไห้\nในชุมชนขุดพบโบราณวัตถุจำนวนมาก เช่น เศษภาชนะดินเผา ลูกปัด พระเครื่องดินเผา และกำไล รวมทั้งเครื่องปั้นดินเผาที่มีลวดลาย เครื่องมือหินอย่างขวานและหินเจียร และเครื่องประดับทำจากหินและกระดูก\nหลักฐานเหล่านี้แสดงว่าบริเวณบ้านเมืองเตาเป็นที่ตั้งถิ่นฐานของผู้คนมาตั้งแต่สมัยก่อนประวัติศาสตร์ และเป็นแหล่งผลิตเหล็กโบราณแห่งหนึ่งในทุ่งกุลาร้องไห้",
   "zh":"孟涛村建在高地上，四周环绕一道宽大的护城河，村庄周围还有古代水塘，据信自古拉龙海铁器时代晚期起就有人在此定居。\n村中出土大量文物，如陶器碎片、珠子、陶制佛牌和手镯，以及带纹饰的陶器、石斧和磨石等石器，还有石制和骨制饰品。\n这些证据表明孟涛一带自史前时代起就有人居住，也是古拉龙海古代炼铁地之一。",
   "en":"Ban Mueang Tao stands on raised ground inside a single wide moat, with ancient ponds around the village, settled since the late Metal Age of the Thung Kula Ronghai culture.\nExcavations in the village have turned up many artefacts: potsherds, beads, terracotta amulets and bangles, decorated pottery, stone tools such as axes and grinding stones, and ornaments of stone and bone.\nTogether they show people have lived here since prehistory, and that this was one of the ancient iron-working sites of Thung Kula Ronghai."}',
 'ข้อเสนอโครงการ FF 2569 หมวด 5','','รอผู้ตรวจสอบชุมชน (A5)','published',4),

('K-005','legend','belief',
 '{"th":"ตำนาน ตอนที่ 1: เซียงสีอพยพจากอัตปือแสนแป","zh":"传说（一）：谢西从阿速坡迁徙而来","en":"Legend, part 1: Xiang Si leaves Attapeu"}',
 '{"th":"ตามตำนานที่ชุมชนเล่าสืบต่อกันมา ผู้สร้างบ้านเมืองเตามาจากเมืองอัตปือแสนแป เมืองเล็ก ๆ ทางทิศตะวันออกเฉียงใต้ของนครจำปาศักดิ์ ประเทศลาว ผู้คนส่วนใหญ่เป็นชาวกุย หรือที่เรียกรวมกันว่าส่วย เป็นคนขยัน อดทน ทำนา ทำไร่ เลี้ยงสัตว์ และชำนาญการจับช้าง\nราว พ.ศ. 2300 เกิดฝนแล้งติดต่อกัน 7–8 ปี ผู้คนจึงอพยพหาที่ทำกินใหม่ โดยมีหัวหน้าเจ็ดคน คือ เซียงฆะ เซียงปูม เซียงสง เซียงสี เซียงพัน เซียงชัย และเซียงขัน แต่ละคนแยกไปตั้งถิ่นฐานที่สังขะ สุรินทร์ เมืองลิง บ้านกุดหวาย บ้านลำดวน บ้านจารพัตร และขุขันธ์ ตามลำดับ\nเซียงสีตั้งถิ่นฐานที่บ้านกุดหวาย (บางตำนานว่าบ้านบุ่งหวาย หรือบ้านเมืองเตา) เป็นผู้มีวิชาอาคมแก่กล้า ชำนาญการจับและเลี้ยงช้าง ทำไร่ทำนา ปลูกถั่วปลูกงา หาปลา และมีฝีมือถลุงเหล็กทำศาสตราวุธ",
   "zh":"据社区代代相传的传说，孟涛的建立者来自老挝占巴塞东南方的小城阿速坡（Attapeu Saen Pae）。当地居民多为奎族，统称“绥族”，勤劳坚忍，种田、饲养牲畜，并擅长捕象。\n约佛历2300年（1757年）连年大旱七八年，人们迁徙寻找新的家园，由七位首领带领：谢卡、谢蒲、谢颂、谢西、谢潘、谢猜和谢坎，分别定居于桑卡、素林、勐林、古外村、兰端村、詹帕村和库坎。\n谢西定居于古外村（另说为汶外村或孟涛村），法力高强，擅长捕象养象，种田、种豆种芝麻、捕鱼，还精于炼铁打造兵器。",
   "en":"As the community tells it, the founder of Mueang Tao came from Attapeu Saen Pae, a small town south-east of Champasak in Laos. Its people were mostly Kui, known together as Suai: hard-working farmers and herders, famed elephant catchers.\nAround 1757 (BE 2300) seven or eight years of drought drove them to seek new land under seven leaders: Xiang Kha, Xiang Pum, Xiang Song, Xiang Si, Xiang Phan, Xiang Chai and Xiang Khan, who settled at Sangkha, Surin, Mueang Ling, Ban Kut Wai, Ban Lamduan, Ban Charaphat and Khukhan.\nXiang Si settled at Ban Kut Wai (some versions say Bung Wai, or Mueang Tao itself). He was a master of spells, skilled at catching and keeping elephants, a farmer and fisher, and a smith who forged weapons."}',
 'หนังสือประวัติเจ้าพ่อศรีนครเตา (แจกในวันฉลองศาล) อ้างใน ข้อเสนอโครงการ FF 2569','','รอผู้ตรวจสอบชุมชน (A5)','published',5),

('K-006','legend','belief',
 '{"th":"ตำนาน ตอนที่ 2: การตามจับช้างเผือก","zh":"传说（二）：追捕白象","en":"Legend, part 2: Catching the white elephants"}',
 '{"th":"ใน พ.ศ. 2302 ปลายสมัยกรุงศรีอยุธยา ช้างเผือกพระที่นั่งแตกโรงหนีไปทางอีสาน เจ้าสองพี่น้อง คือพระยาจักรี (นายทองด้วง) และพระยาสุรสีห์ (นายบุญมา) นำขบวนราว 30 คนออกติดตาม ผ่านเมืองพิมายจนถึงทุ่งกุลาร้องไห้ และได้พบเซียงสี\nเพราะต่างภาษากันจึงสื่อสารกันลำบาก ฝ่ายหนึ่งไม่เข้าใจคำว่า \"ด่อน\" อีกฝ่ายไม่เข้าใจคำว่า \"เผือก\" สุดท้ายเจ้าสองพี่น้องได้ผูกเสี่ยวกับเซียงสีและขอให้ช่วยตามช้าง\nคณะติดตามผ่านทุ่งกุลาร้องไห้ บ้านห้วยแก้ว (อำเภอรัตนบุรีในปัจจุบัน) จนพบพญาช้างเผือก 2 เชือกเล่นน้ำอยู่กับโขลงช้างป่าที่หนองบัวขี้ตะโหลก เซียงสีเสกก้อนดิน 8 ก้อนโยนไปแปดทิศ ร่ายคาถาเชิญพญาช้างกลับกรุงศรีอยุธยา ช้างป่าแตกตื่นหนีไป เหลือเพียงช้างเผือกที่ขึ้นฝั่งหมอบลงต่อหน้า สถานที่นั้นจึงได้ชื่อว่า \"ห้วยทัพทัน\"",
   "zh":"佛历2302年（1759年），阿瑜陀耶王朝末期，御用白象挣脱象厩向东北逃去。兄弟二人——帕耶却克里（通銮）与帕耶素拉西（汶玛）率约30人追寻，经披迈抵达古拉龙海平原，遇见了谢西。\n因语言不通，双方难以沟通：一方不懂“端”（白），另一方不懂“培”（白色）。最终兄弟二人与谢西结为“西奥”（结拜兄弟），请他协助追象。\n众人穿过古拉龙海、怀乔村（今叻达那武里县），在农布基达洛发现两头白象正与五六十头野象戏水。谢西念咒将八块土块抛向八方，恭请白象返回阿瑜陀耶。野象受惊四散，只剩白象上岸跪伏在他面前，此地因而得名“怀塔坦”。",
   "en":"In 1759 (BE 2302), late in the Ayutthaya era, the royal white elephant broke loose and fled north-east. Two brothers, Phraya Chakri (Thongduang) and Phraya Surasi (Bunma), followed with about thirty men through Phimai to Thung Kula Ronghai, where they met Xiang Si.\nThey spoke different languages: one side did not know the word don, the other did not know phueak, both meaning white. In the end the brothers became Xiang Si sworn friends (siao) and asked his help.\nThe party crossed Thung Kula and Ban Huai Kaeo (today Rattanaburi) and found two white elephants bathing with a wild herd at Nong Bua Khi Talok. Xiang Si charmed eight clods of earth, threw them to the eight directions and invited the elephants home. The herd scattered; the white elephants came ashore and knelt before him. The place became Huai Thap Than."}',
 'หนังสือประวัติเจ้าพ่อศรีนครเตา อ้างใน ข้อเสนอโครงการ FF 2569','','รอผู้ตรวจสอบชุมชน (A5)','published',6),

('K-007','legend','belief',
 '{"th":"ตำนาน ตอนที่ 3: ดาบนวโลหะกับบ่อแท่งเหล็กแทงทอง","zh":"传说（三）：九金之剑与“插铁入金井”","en":"Legend, part 3: The nine-metal sword"}',
 '{"th":"เจ้าสองพี่น้องถูกใจดาบของเซียงสีจึงขอไว้ เซียงสีจึงบอกว่าดาบเก่าแล้ว ถ้าอยากได้ให้หาเหล็ก 7 แหย่งช้างมา จะตีให้ใหม่ เมื่อได้เหล็กแล้ว ลูกทั้งหกของเซียงสีช่วยกันทำเตาถลุงเหล็ก ต่อมาบริเวณนั้นกลายเป็นแอ่งน้ำที่ชาวบ้านเรียกว่า \"หนองเตาเหล็ก\"\nเซียงสีตีดาบด้วย \"นวโลหะ\" คือโลหะเก้าชนิด ได้แก่ เหล็กไหล เหล็กปลีก เหล็กกล้า ทองเหลือง ทองแดง นาก เงิน ทองคำ และทองขาว เริ่มตั้งแต่ขึ้น 1 ค่ำ เดือน 9 ปีมะเส็ง จนแล้วเสร็จในวันเพ็ญเดือนสาม พ.ศ. 2304 รวม 6 เดือน 15 วัน เมื่อนำไปลองฟันไม้ที่ดงมูลครั้งเดียวขาดแสนต้น ที่นั่นจึงเรียกว่า \"ดงสะตาว\" หรือ \"ดงแสนตอ\"\nเมื่อนำดาบไปถวายที่กรุงศรีอยุธยา เจ้าสองพี่น้องไม่พอใจที่เหล็กมากมายได้ดาบเพียงเล่มเดียว ขณะเซียงสีลุกขึ้นลากลับ ท้องพระโรงสั่นสะเทือน ดาบกระเด็นกระดอน เซียงสีจึงคว้าโยนออกไปนอกท้องพระโรง ปลายดาบปักลงดินจนแยกเป็นบ่อน้ำใหญ่ เรียกว่า \"บ่อแท่งเหล็กแทงทอง\" หรือ \"สระกัดเหล็กกัดทอง\" และดาบก็จมหายไป",
   "zh":"兄弟二人看中谢西的佩剑想要索取。谢西说这把剑已经旧了，若真想要，就带来七驮象的铁，他会重新打造。铁运来后，谢西的六个孩子一起建造炼铁炉，那里后来成为一片水洼，村民称之为“铁炉塘”。\n谢西用“九金”（九种金属：陨铁、熟铁、钢、黄铜、紫铜、红金、银、金、白金）铸剑，自蛇年九月初一开工，至佛历2304年三月月圆日完成，历时六个月十五天。在“蒙林”试剑，一挥便断十万棵树，此地因此称为“萨涛林”或“十万桩林”。\n剑献至阿瑜陀耶时，兄弟二人不满如此多的铁只铸成一把剑。谢西起身告辞时，殿堂震动，剑四处弹跳，他抓起剑抛出殿外，剑尖插入地中，大地裂开成为一口大井，称为“插铁入金井”，剑也沉没不见。",
   "en":"The two brothers admired Xiang Si sword and asked for it. He said it was old; if they truly wanted one, they should bring iron by the elephant-load, seven loads, and he would forge a new one. When the iron came, his six children built a smelting furnace on a spot that later filled with water, the pond villagers still call Nong Tao Lek.\nHe forged the blade from nawa loha, nine metals: meteoric iron, wrought iron, steel, brass, copper, nak, silver, gold and white gold, from the first waxing day of the ninth month to the full moon of the third month of 1761, six months and fifteen days. Tested in the Dong Mun forest, one stroke felled a hundred thousand trees, so the place is called Dong Sa Tao.\nAt Ayutthaya the brothers were unhappy that so much iron had made a single sword. As Xiang Si rose to leave, the hall shook and the sword leapt about; he flung it outside, its point split the earth into a great well, Bo Thaeng Lek Thaeng Thong, and the sword sank from sight."}',
 'หนังสือประวัติเจ้าพ่อศรีนครเตา อ้างใน ข้อเสนอโครงการ FF 2569','','รอผู้ตรวจสอบชุมชน (A5)','published',7),

('K-008','legend','belief',
 '{"th":"ตำนาน ตอนที่ 4: จากเซียงสีสู่ \"เจ้าพ่อศรีนครเตา\"","zh":"传说（四）：从谢西到“西那空涛神”","en":"Legend, part 4: From Xiang Si to Chao Pho Si Nakhon Tao"}',
 '{"th":"เหตุการณ์เรื่องดาบทำให้เจ้าสองพี่น้องเห็นความซื่อสัตย์ของเซียงสี จึงนำทั้งหกคนเข้าเฝ้า และได้รับพระราชทานบรรดาศักดิ์ เซียงสีเป็น \"หลวงศรีนครเตาท้าวเธอ\" และบ้านกุดหวายได้ยกฐานะเป็นบ้านเมืองเตา ท่านใช้เวลา 3 ปีขุดคูน้ำล้อมรอบเมือง ซึ่งปัจจุบันเรียกว่าลำน้ำเตา\nท่านรับราชการด้วยความจงรักภักดี ช่วยราชการสงครามหลายครั้ง ดูแลทุกข์สุขของราษฎร จนได้เลื่อนเป็น \"พระศรีนครเตาท้าวเธอ\" ที่คนทั่วไปเรียกว่า \"เจ้าพ่อศรีนครเตา\"\nในสมัยรัชกาลที่ 1 ใบบอกราชการที่ส่งไปยังเมืองอุบลตกหล่นที่เมืองเตา ทำให้ทัพเมืองอุบลไปช่วยเชียงใหม่ไม่ทัน จึงมีรับสั่งให้ลงโทษ ท่านจึงหนีไปบวชที่บ้านไพรขลา เปลี่ยนชื่อเป็น \"เมื่อย\" ชาวบ้านมาเล่าเรียนวิชาจากท่านมากมาย จนเป็นพระอุปัชฌาย์นาม \"ฌาย์เมื่อย\" ตามบันทึกคำบอกเล่าของคุณปู่เสน ผาจีบ ท่านมรณภาพที่วัดไพรขลาเมื่อ พ.ศ. 2338 อายุ 86 ปี\nต่อมาลูกหลานบางส่วนย้ายไปตั้งเมืองห้วยแก้ว (รัตนบุรี) และบุตรของท่านได้ครองเมืองต่าง ๆ ในแถบทุ่งกุลา",
   "zh":"剑的事件让兄弟二人看到了谢西的忠诚，于是带六人觐见国王并获封爵位：谢西受封为“銮·西那空涛陶特”，古外村升格为孟涛村。他花了三年时间开挖环城水道，即今日的涛河。\n他忠心效力，多次参战，体恤百姓，后晋升为“帕·西那空涛陶特”，百姓尊称他为“西那空涛神”（Chao Pho Si Nakhon Tao）。\n拉玛一世时期，送往乌汶的公文在孟涛遗失，致使乌汶军未能及时援助清迈，国王下令惩处。他逃往派卡拉村出家，改名“蒙伊”。许多村民前来向他学习，他成为戒师“差·蒙伊”。据塞恩·帕集老人的口述记录，他于佛历2338年（1795年）在派卡拉寺圆寂，享年86岁。\n此后部分后人迁往怀乔（叻达那武里）建城，他的子女分别治理古拉龙海一带的城镇。",
   "en":"The sword episode convinced the brothers of Xiang Si honesty, and they presented all six to the king. Xiang Si was made Luang Si Nakhon Tao Thao Thoe and Ban Kut Wai was raised to Ban Mueang Tao. He spent three years digging the moat around the town, today called the Tao stream.\nHe served loyally, fought in several campaigns and cared for his people, and was raised to Phra Si Nakhon Tao Thao Thoe, known to all as Chao Pho Si Nakhon Tao.\nUnder King Rama I a royal dispatch to Ubon was lost at Mueang Tao, so the Ubon army failed to reach Chiang Mai in time, and punishment was ordered. He fled to Ban Phrai Khla, was ordained under the name Muai, taught many villagers and became the preceptor Chai Muai. By the account of Grandfather Sen Phachip he died at Wat Phrai Khla in 1795 (BE 2338), aged 86.\nSome descendants later founded Huai Kaeo (Rattanaburi), and his children governed towns across the Thung Kula plain."}',
 'หนังสือประวัติเจ้าพ่อศรีนครเตา; บันทึกปากคำคุณปู่เสน ผาจีบ (อ้างใน ข้อเสนอโครงการ FF 2569)','','รอผู้ตรวจสอบชุมชน (A5)','published',8),

('K-009','belief','belief',
 '{"th":"ความเชื่อเรื่องผีในชุมชนบ้านเมืองเตา","zh":"孟涛村的神灵信仰","en":"Spirit beliefs in Ban Mueang Tao"}',
 '{"th":"ความเชื่อเรื่องผีเป็นส่วนหนึ่งของวิถีชีวิตชาวบ้านเมืองเตามายาวนาน โดยเฉพาะความเชื่อเรื่องเจ้าพ่อศรีนครเตา ซึ่งชาวบ้านนับถือเป็นผีปู่ตาและบรรพบุรุษผู้สร้างบ้านแปงเมือง ผูกพันกับชุมชนมาหลายชั่วอายุคน\nนอกจากนี้ยังมี ผีบ้านผีเรือน ที่เชื่อว่าคอยดูแลทุกบ้าน จึงมีการเซ่นไหว้เป็นประจำ ผีป่า ที่ทำให้ชาวบ้านระมัดระวังและเคารพธรรมชาติเมื่อเข้าป่า และผีน้ำ ที่ทำให้ชาวบ้านขอขมาก่อนลงเล่นน้ำหรือใช้น้ำ\nความเชื่อนี้มีผลต่อการทำนา เช่น การทำพิธีขอฝนและขอขมาสิ่งศักดิ์สิทธิ์ก่อนฤดูทำนา การเลี้ยงสัตว์ และการรักษาโรคด้วยหมอพื้นบ้านและสมุนไพรควบคู่กับการแพทย์สมัยใหม่\nการนับถือผีอยู่ร่วมกับความเชื่อแบบพราหมณ์และพุทธศาสนา รวมถึงฮีตสิบสองคองสิบสี่ ซึ่งช่วยให้ชาวบ้านวางแผนชีวิตให้สอดคล้องกับฤดูกาล",
   "zh":"神灵信仰长期以来是孟涛村民生活的一部分，尤其是对西那空涛神的信仰。村民将他奉为“祖灵”（Phi Pu Ta），即建村立城的祖先，与社区相连已有许多代人。\n此外还有守护每户人家的家神，村民定期祭拜；林神，使村民入林时谨慎并敬畏自然；以及水神，村民下水或用水前会先请求宽恕。\n这些信仰影响着农耕，例如插秧季节前举行祈雨仪式并向神灵告罪；也影响着饲养牲畜，以及与现代医学并行的民间医师和草药治疗。\n神灵信仰与婆罗门教、佛教以及“十二习俗十四规范”（Heet Sip Song Khong Sip Si）并存，帮助村民按季节安排生活。",
   "en":"Belief in spirits has long shaped life in Ban Mueang Tao, above all belief in Chao Pho Si Nakhon Tao, honoured as the phi pu ta, the ancestor spirit who founded the town, and bound to the village for generations.\nVillagers also keep house spirits, with regular offerings in every home; forest spirits, which make people careful and respectful in the woods; and water spirits, whose pardon is asked before swimming in or using a pond.\nThese beliefs shape farming, with a rain-calling and apology rite before planting, as well as livestock keeping and healing, where folk healers and herbs sit alongside modern medicine.\nSpirit worship lives together with Brahmin and Buddhist practice and the Isan calendar of Heet Sip Song Khong Sip Si, helping villagers plan life around the seasons."}',
 'ข้อเสนอโครงการ FF 2569 หมวด 5 (ความสำคัญของชุมชน ข้อ 1)','','รอผู้ตรวจสอบชุมชน (A5)','published',9),

('K-010','belief','belief',
 '{"th":"บทบาทของความเชื่อเจ้าพ่อศรีนครเตา","zh":"西那空涛神信仰的作用","en":"What the belief in Chao Pho does for people"}',
 '{"th":"การศึกษาพิธีกรรมพบว่าความเชื่อเจ้าพ่อศรีนครเตามีบทบาทสองระดับ\nระดับปัจเจกบุคคล เป็นที่พึ่งทางใจในยามวิกฤติหรือยากลำบาก เช่น เมื่อไม่สบายใจ เจ็บป่วย เดินทางไกล หรือสอบแข่งขัน ช่วยให้ผู้คนมีความหวังและรู้สึกมั่นคงขึ้น\nระดับชุมชน ช่วยรับรองสถานภาพใหม่ของบุคคล ยืนยันความศักดิ์สิทธิ์ของพิธีกรรม สนับสนุนวงจรการเกษตรและความอุดมสมบูรณ์ และยังช่วยควบคุมสังคม สร้างความสามัคคี รักษาโครงสร้างของชุมชนให้มั่นคง\nชาวบ้านเชื่อว่าเจ้าพ่อเป็นผู้ปกปักรักษาชุมชนและให้ความอุดมสมบูรณ์ทางการเกษตร จึงจัดงานบวงสรวงและขอพรเป็นประจำทุกปี",
   "zh":"仪式研究发现，西那空涛神信仰在两个层面发挥作用。\n个人层面：在危机或困难时成为精神依靠，例如心中不安、生病、远行或考试时，给人希望和安全感。\n社区层面：认可个人的新身份，确认仪式的神圣性，支持农耕周期与丰收，同时有助于维护社会秩序、凝聚团结、稳定社区结构。\n村民相信神灵守护社区、赐予农业丰收，因此每年都举行祭祀与祈福。",
   "en":"Studies of the rites find the belief works at two levels.\nFor the individual it is a refuge in hard times: worry, illness, a long journey, an examination. It gives hope and a sense of security.\nFor the community it confirms a person new status, affirms the sanctity of rites, supports the farming cycle and fertility, and keeps social order, unity and the village structure steady.\nBelieving that Chao Pho protects the village and its harvests, villagers hold an offering and blessing ceremony every year."}',
 'อมาวสี เถียรถาวร (2532); ข้อเสนอโครงการ FF 2569 หมวด 5–6','','รอผู้ตรวจสอบชุมชน (A5)','published',10),

('K-011','belief','belief',
 '{"th":"ความศรัทธาในทุ่งกุลาร้องไห้","zh":"古拉龙海平原上的信仰","en":"Faith across Thung Kula Ronghai"}',
 '{"th":"ความเชื่อเจ้าพ่อศรีนครเตาพบได้ทั่วทุ่งกุลาร้องไห้ ทั้งชุมชนตอนกลางของทุ่ง เช่น บ้านเมืองบัว และบ้านกู่กาสิงห์ อำเภอเกษตรวิสัย จังหวัดร้อยเอ็ด ชุมชนชายขอบทุ่งด้านใต้ เช่น บ้านไพรขลา อำเภอชุมพลบุรี จังหวัดสุรินทร์ และชุมชนชายขอบด้านตะวันตก คือบ้านเมืองเตาและบ้านเมืองเสือ อำเภอพยัคฆภูมิพิสัย จังหวัดมหาสารคาม\nแต่ละชุมชนมีรูปแบบความเชื่อต่างกันตามบริบทและภูมินิเวศของตน\nบ้านเมืองเตามีความศรัทธาเข้มข้นเป็นพิเศษ เพราะเชื่อว่าเจ้าพ่อคือบรรพบุรุษผู้สร้างบ้านแปลงเมือง และตั้งชื่อชุมชนตามชื่อท่าน จึงเกิดเป็นความทรงจำร่วมของชุมชน แสดงออกผ่านงานเฉลิมฉลองในเดือน 4 ตามจันทรคติของทุกปี",
   "zh":"西那空涛神信仰遍布整个古拉龙海平原：平原中部的社区，如黎逸府格塞威赛县的孟布村和古卡辛村；南缘的社区，如素林府春蓬武里县的派卡拉村；以及西缘的社区，即玛哈沙拉堪府帕亚卡普米披赛县的孟涛村和孟苏村。\n各社区的信仰形式因其环境与生态而有所不同。\n孟涛村的信仰尤为浓厚，因为村民相信神灵是建村立城的祖先，村名也随其名号而来，从而形成社区的集体记忆，并通过每年农历四月的庆典表现出来。",
   "en":"Belief in Chao Pho Si Nakhon Tao runs across Thung Kula Ronghai: in villages at its heart such as Ban Mueang Bua and Ban Ku Ka Sing (Kaset Wisai, Roi Et), on its southern edge at Ban Phrai Khla (Chumphon Buri, Surin), and on its western edge at Ban Mueang Tao and Ban Mueang Suea (Phayakkhaphum Phisai, Maha Sarakham).\nEach community holds the belief in its own way, shaped by its setting.\nIt is strongest in Ban Mueang Tao, which sees Chao Pho as the ancestor who founded the town and named the village after him. That shared memory is renewed at the festival held every year in the fourth lunar month."}',
 'ข้อเสนอโครงการ FF 2569 หมวด 6; พิเนตร ดาวเรือง (2563)','','รอผู้ตรวจสอบชุมชน (A5)','published',11),

('K-012','community','history',
 '{"th":"คำขวัญตำบลเมืองเตา","zh":"孟涛乡标语","en":"The Mueang Tao subdistrict motto"}',
 '{"th":"\"ถิ่นข้าวหอมมะลิหนึ่งศูนย์ห้า ทุ่งกุลาสดใส งามวิไลไหมมัดหมี่ รู้กันดีหมวกจักสาน สืบตำนานงานบั้งไฟ สมุนไพรเลื่องลือ นับถือเจ้าพ่อศรีนครเตา\"\nคำขวัญนี้มาจากการประกวดคำขวัญประจำตำบลที่องค์การบริหารส่วนตำบลเมืองเตาจัดขึ้นเมื่อ พ.ศ. 2542 และยังใช้มาจนถึงปัจจุบัน สะท้อนของดีของตำบล ได้แก่ ข้าวหอมมะลิ 105 ทุ่งกุลาร้องไห้ ผ้าไหมมัดหมี่ หมวกจักสาน งานบุญบั้งไฟ สมุนไพร และความศรัทธาต่อเจ้าพ่อศรีนครเตา",
   "zh":"“茉莉香米105之乡，古拉平原明亮，扎染丝绸秀丽，编织帽远近闻名，传承火箭节传统，草药声名远扬，敬奉西那空涛神。”\n这条标语来自孟涛乡行政机构于佛历2542年（1999年）举办的乡标语征集活动，沿用至今，体现了本乡的特色：茉莉香米105、古拉龙海平原、扎染丝绸、编织帽、火箭节、草药以及对西那空涛神的信仰。",
   "en":"Home of Hom Mali 105 jasmine rice, the bright Kula plain, lovely mudmee silk, well-known woven hats, the Bun Bang Fai rocket festival, famed herbs, and faith in Chao Pho Si Nakhon Tao.\nThe motto won a contest held by the Mueang Tao Subdistrict Administrative Organisation in 1999 (BE 2542) and is still in use. It lists what the subdistrict is proud of: jasmine rice, the Thung Kula plain, mudmee silk, woven hats, the rocket festival, herbs and the shrine."}',
 'จิรายุทธ แสงมณี (ปราชญ์ชุมชน) สัมภาษณ์ 16 เม.ย. 2567 โดย นัชชา อู่เงิน','จิรายุทธ แสงมณี','รอผู้ตรวจสอบชุมชน (A5)','published',12),

('K-013','community','community',
 '{"th":"วิถีชีวิตและสังคมของชาวบ้านเมืองเตา","zh":"孟涛村民的生活与社会","en":"Life and society in Ban Mueang Tao"}',
 '{"th":"ชุมชนบ้านเมืองเตายังคงวิถีชีวิตและวัฒนธรรมแบบชนบท พื้นที่รอบชุมชนเป็นนาข้าว ชาวบ้านส่วนใหญ่ทำเกษตรกรรมแบบยั่งยืน เช่น ปลูกพืชหมุนเวียน ใช้ปุ๋ยธรรมชาติ และร่วมกันอนุรักษ์พื้นที่ชุ่มน้ำของทุ่งกุลาร้องไห้\nการตั้งถิ่นฐานหลังแต่งงานส่วนใหญ่อยู่ฝ่ายหญิง ลูกหลานจึงใกล้ชิดญาติฝ่ายแม่ นามสกุลเก่าแก่ที่ใช้แพร่หลาย ได้แก่ ไชยสงคราม บุราณรมณ์ แก้วประเสริฐ อาษากิจ ประทุมทอง พลอามาตย์ และพงษ์หนองพอก\nชุมชนมีกิจกรรมส่งเสริมการเรียนรู้และอาชีพ เช่น การทอผ้า และการทำอาหารพื้นเมือง ชาวบ้านช่วยเหลือเกื้อกูลกันและให้ความสำคัญกับการสืบทอดประเพณี\nปราชญ์ชุมชนที่ให้ข้อมูลแก่โครงการ ได้แก่ นายจิรายุทธ แสงมณี นายสวานิตย์ ศิริบูรณ์ และนายทรงเดช บัวไชยา",
   "zh":"孟涛村仍保持着乡村的生活方式与文化。村庄周围是稻田，大多数村民从事可持续农业，例如轮作、使用天然肥料，并共同保护古拉龙海的湿地。\n婚后大多随女方居住，因此子女与母系亲属更为亲近。常见的老姓氏有：猜颂克兰、武兰隆、乔巴瑟、阿萨吉、巴吞通、蓬阿玛和蓬农波。\n社区开展织布、制作地方菜等学习与职业活动。村民互帮互助，重视传统的传承。\n为本项目提供信息的社区学者有：吉拉育·桑玛尼、萨瓦尼·西里汶和宋德·布猜亚。",
   "en":"Ban Mueang Tao keeps a rural way of life among its rice fields. Most villagers farm sustainably, rotating crops, using natural fertiliser and caring together for the Thung Kula wetlands.\nAfter marriage couples usually live with the wife family, so children grow up close to their mother relatives. Old, widespread family names include Chaisongkhram, Buranarom, Kaeoprasoet, Asakit, Prathumthong, Phon-amat and Phongnongphok.\nThe village runs learning and livelihood activities such as weaving and local cooking, helps one another and cares about passing traditions on.\nCommunity scholars who informed the project are Mr Chirayut Saengmani, Mr Sawanit Siriboon and Mr Songdet Buachaiya."}',
 'ข้อเสนอโครงการ FF 2569 หมวด 5 (ความสำคัญของชุมชน ข้อ 2–3)','','รอผู้ตรวจสอบชุมชน (A5)','published',13),

('K-014','ritual','ritual',
 '{"th":"แนวทางการสักการะเจ้าพ่อศรีนครเตา","zh":"参拜西那空涛神的建议","en":"How to pay respect at the shrine"}',
 '{"th":"แนวทางทั่วไปสำหรับผู้มาเยือน ขั้นตอนเฉพาะของชุมชนรอการยืนยันจากผู้รู้ในชุมชน\n1. แต่งกายสุภาพ\n2. เตรียมดอกไม้ ธูป เทียน\n3. จุดธูปเทียน กล่าวคำบูชา\n4. อธิษฐานด้วยความสงบ ไม่ส่งเสียงดัง\n5. เมื่อสมหวังจึงกลับมาแก้บนตามที่ได้ตั้งจิตไว้",
   "zh":"以下为访客的一般建议，社区特有的步骤待社区长老确认。\n1. 穿着得体\n2. 准备鲜花、香和蜡烛\n3. 点香烛，诵读祭词\n4. 安静地祈祷，不要喧哗\n5. 愿望实现后，按许下的承诺回来还愿",
   "en":"General guidance for visitors; the community own steps await confirmation by local elders.\n1. Dress modestly\n2. Bring flowers, incense and candles\n3. Light the incense and candles and say the words of offering\n4. Pray quietly\n5. If your wish comes true, return to give thanks as you promised"}',
 'แนวทางทั่วไป รอข้อมูลพิธีกรรมจากชุมชน (A4)','','รอผู้ตรวจสอบชุมชน (A5)','published',14),

('K-015','ritual','ritual',
 '{"th":"พิธีขอฝนก่อนฤดูทำนา","zh":"插秧季前的祈雨仪式","en":"Calling the rain before planting"}',
 '{"th":"ก่อนเริ่มฤดูทำนา ชาวบ้านเมืองเตาทำพิธีขอฝนและขอขมาต่อสิ่งศักดิ์สิทธิ์ เพื่อให้ฝนตกต้องตามฤดูกาลและได้ผลผลิตที่ดี พิธีนี้สะท้อนบทบาทของเจ้าพ่อศรีนครเตาในฐานะผู้ปกปักรักษาชุมชนและผู้ให้ความอุดมสมบูรณ์ทางการเกษตร",
   "zh":"在插秧季节开始前，孟涛村民举行祈雨仪式并向神灵告罪，祈求风调雨顺、五谷丰登。这一仪式体现了西那空涛神作为社区守护者和农业丰收赐予者的角色。",
   "en":"Before the rice season begins, villagers hold a rite to call the rain and ask the sacred powers for pardon, so the rain comes on time and the harvest is good. It reflects Chao Pho role as guardian of the village and giver of plenty."}',
 'ข้อเสนอโครงการ FF 2569 หมวด 5','','รอผู้ตรวจสอบชุมชน (A5)','published',15);

-- ---------------------------------------------------------------- places
-- Only the village coordinate (15.5116, 103.3866) is given in the source; site positions are approximate.
INSERT INTO places(code, category, name, description, lat, lon, open_hours, visit_minutes, tags, source) VALUES
('P-001','sacred','{"th":"ศาลเจ้าพ่อศรีนครเตา","zh":"西那空涛神祠","en":"Chao Pho Si Nakhon Tao Shrine"}',
 '{"th":"ศูนย์รวมศรัทธาของชาวบ้านเมืองเตาและผู้คนในทุ่งกุลาร้องไห้ ชาวบ้านนับถือเจ้าพ่อศรีนครเตาเป็นผีปู่ตาและบรรพบุรุษผู้สร้างบ้านแปลงเมือง ทุกปีมีงานบวงสรวงและเฉลิมฉลองในเดือน 4 ตามจันทรคติ\nผู้คนมาขอพรในยามไม่สบายใจ เจ็บป่วย เดินทางไกล หรือก่อนสอบแข่งขัน",
   "zh":"孟涛村民和古拉龙海民众的信仰中心。村民将西那空涛神奉为祖灵和建村立城的祖先，每年农历四月举行祭祀与庆典。\n人们在心中不安、生病、远行或考试前来此祈福。",
   "en":"The heart of faith for Ban Mueang Tao and the Thung Kula plain. Villagers honour Chao Pho as their ancestor spirit and founder, with an offering festival each fourth lunar month.\nPeople come to pray when worried, ill, before a long journey or an examination."}',
 15.5118,103.3868,'06:00–18:00',45,'{shrine,worship,legend}','ข้อเสนอโครงการ FF 2569; พิกัดโดยประมาณ'),
('P-002','history','{"th":"แหล่งโบราณคดีบ้านเมืองเตา","zh":"孟涛村考古遗址","en":"Ban Mueang Tao archaeological site"}',
 '{"th":"ชุมชนโบราณบนที่ดอน มีคูน้ำคันดินขนาดใหญ่ล้อมรอบชั้นเดียว อายุตั้งแต่ก่อนประวัติศาสตร์จนถึงสมัยวัฒนธรรมเขมร ตามการสำรวจของกรมศิลปากร (2529–2546)\nพบเศษภาชนะดินเผา ลูกปัด พระเครื่องดินเผา กำไล เครื่องมือหิน และเครื่องประดับจากหินและกระดูก",
   "zh":"位于高地上的古代聚落，四周环绕一道宽大的护城河与土墙，年代从史前延续到高棉文化时期（泰国艺术厅1986–2003年调查）。\n出土陶器碎片、珠子、陶制佛牌、手镯、石器以及石制和骨制饰品。",
   "en":"An ancient settlement on raised ground inside a single wide moat and earthen bank, from prehistory to the Khmer period (Fine Arts Department survey, 1986–2003).\nFinds include potsherds, beads, terracotta amulets, bangles, stone tools and ornaments of stone and bone."}',
 15.5130,103.3885,'ตลอดวัน',40,'{archaeology,moat,metal_age}','กรมศิลปากร (อ้างใน ข้อเสนอโครงการ FF 2569); พิกัดโดยประมาณ'),
('P-003','history','{"th":"หนองเตาเหล็ก","zh":"铁炉塘","en":"Nong Tao Lek, the iron-furnace pond"}',
 '{"th":"แหล่งถลุงเหล็กสมัยก่อนประวัติศาสตร์และที่มาของชื่อ \"เมืองเตา\" อยู่ทางทิศตะวันตกเฉียงใต้ของหมู่บ้าน บริเวณหลังโรงเรียนเมืองเตาวิทยาคม ยังพบตะกรันเหล็กและศิลาแลงอยู่มาก\nตามตำนาน ที่นี่คือที่ลูกทั้งหกของเซียงสีสร้างเตาถลุงเหล็กเพื่อตีดาบนวโลหะ",
   "zh":"史前炼铁遗址，也是“孟涛”之名的由来。位于村庄西南方、孟涛威塔亚空中学后方，仍可见大量铁渣与红土石。\n据传说，谢西的六个孩子在此建炉炼铁，打造九金之剑。",
   "en":"A prehistoric iron-smelting site that gave Mueang Tao its name, south-west of the village behind Mueang Tao Witthayakhom School, where slag and laterite are still plentiful.\nIn the legend, Xiang Si six children built the furnace here to forge the nine-metal sword."}',
 15.5085,103.3832,'ตลอดวัน',30,'{archaeology,iron,legend}','ข้อเสนอโครงการ FF 2569; พิกัดโดยประมาณ'),
('P-004','culture','{"th":"ลำน้ำเตา","zh":"涛河","en":"The Tao stream"}',
 '{"th":"ลำน้ำสายเล็กที่ไหลผ่านห่างจากชุมชนราว 2 กิโลเมตร เป็นแหล่งน้ำสำคัญเพื่อการอุปโภคบริโภคและการเกษตร ตามตำนาน หลวงศรีนครเตาใช้เวลา 3 ปีขุดคูน้ำล้อมรอบเมือง ซึ่งปัจจุบันเรียกว่าลำน้ำเตา\nชาวบ้านเชื่อเรื่องผีน้ำ จึงขอขมาก่อนลงเล่นน้ำหรือใช้น้ำ",
   "zh":"距村约2公里的小河，是生活与农业的重要水源。据传说，銮·西那空涛花了三年时间开挖环城水道，即今日的涛河。\n村民相信水神，下水或用水前会先请求宽恕。",
   "en":"A small stream about 2 km from the village and its main water for homes and fields. Legend says Luang Si Nakhon Tao spent three years digging the moat around the town, now called the Tao stream.\nVillagers believe in water spirits and ask their pardon before swimming or drawing water."}',
 15.5240,103.3945,'ตลอดวัน',30,'{nature,legend,wetland}','ข้อเสนอโครงการ FF 2569; พิกัดโดยประมาณ');

INSERT INTO routes(name, description, place_ids, minutes) VALUES
('{"th":"เส้นทางตำนานเมืองเตา ครึ่งวัน","zh":"孟涛传说半日路线","en":"Half-day Mueang Tao legend route"}',
 '{"th":"สักการะเจ้าพ่อ ชมคูน้ำคันดินเมืองโบราณ เยือนหนองเตาเหล็ก และลำน้ำเตา","zh":"参拜神祠、参观古城护城河、探访铁炉塘与涛河","en":"The shrine, the ancient moat, the iron-furnace pond and the Tao stream"}',
 ARRAY(SELECT id FROM places WHERE code IN ('P-001','P-002','P-003','P-004') ORDER BY code), 180);

-- ---------------------------------------------------------------- events
-- Lunar-calendar dates move every year; the community announces the exact day (survey item A4).
INSERT INTO events(title, description, preparation, starts_at, ends_at, place_id, source) VALUES
('{"th":"งานบวงสรวงและเฉลิมฉลองเจ้าพ่อศรีนครเตา","zh":"西那空涛神祭祀庆典","en":"Annual offering festival of Chao Pho Si Nakhon Tao"}',
 '{"th":"งานประจำปีในเดือน 4 ตามจันทรคติ ชาวบ้านเมืองเตาและผู้ศรัทธาจากทั่วทุ่งกุลาร้องไห้ร่วมบวงสรวงและขอพร วันที่แสดงเป็นวันโดยประมาณ วันที่แน่นอนประกาศโดยชุมชน","zh":"每年农历四月举行，孟涛村民和古拉龙海各地信众共同祭祀祈福。所示日期为估计日期，确切日期由社区公布。","en":"Held each fourth lunar month, when villagers and believers from across Thung Kula make offerings and ask for blessings. The date shown is an estimate; the community announces the exact day."}',
 '{"th":"แต่งกายสุภาพ เตรียมดอกไม้ ธูป เทียน","zh":"穿着得体，准备鲜花香烛","en":"Dress modestly; bring flowers, incense and candles"}',
 '2027-03-06 08:00+07','2027-03-07 17:00+07',(SELECT id FROM places WHERE code='P-001'),'ข้อเสนอโครงการ FF 2569 หมวด 6'),
('{"th":"พิธีขอฝนก่อนฤดูทำนา","zh":"插秧季前祈雨仪式","en":"Rain-calling rite before planting"}',
 '{"th":"ชาวบ้านทำพิธีขอฝนและขอขมาสิ่งศักดิ์สิทธิ์ เพื่อให้ฝนตกต้องตามฤดูกาล วันที่แสดงเป็นวันโดยประมาณ","zh":"村民举行祈雨仪式并向神灵告罪，祈求风调雨顺。所示日期为估计日期。","en":"Villagers call the rain and ask the sacred powers for pardon before planting. The date shown is an estimate."}',
 '{"th":"แต่งกายสุภาพ เคารพพิธีกรรมของชุมชน","zh":"穿着得体，尊重社区仪式","en":"Dress modestly and respect the community rite"}',
 '2027-05-15 08:00+07','2027-05-15 12:00+07',(SELECT id FROM places WHERE code='P-001'),'ข้อเสนอโครงการ FF 2569 หมวด 5'),
('{"th":"งานบุญบั้งไฟ","zh":"火箭节","en":"Bun Bang Fai rocket festival"}',
 '{"th":"ประเพณีบุญบั้งไฟที่ปรากฏในคำขวัญตำบลเมืองเตา \"สืบตำนานงานบั้งไฟ\" วันที่แสดงเป็นวันโดยประมาณ","zh":"孟涛乡标语中“传承火箭节传统”所指的火箭节。所示日期为估计日期。","en":"The rocket festival named in the subdistrict motto. The date shown is an estimate."}',
 '{"th":"ระวังจุดปล่อยบั้งไฟ ปฏิบัติตามคำแนะนำของผู้จัดงาน","zh":"注意火箭发射点，遵守主办方指引","en":"Keep clear of the launch area and follow the organisers"}',
 '2027-06-05 09:00+07','2027-06-06 18:00+07',NULL,'คำขวัญตำบลเมืองเตา (2542)');

-- ---------------------------------------------------------------- faq
INSERT INTO faq(category, question, answer, doc_code) VALUES
('history','{"th":"ศาลเจ้าพ่อศรีนครเตาอยู่ที่ไหน","zh":"西那空涛神祠在哪里？","en":"Where is the shrine?"}',
 '{"th":"บ้านเมืองเตา ตำบลเมืองเตา อำเภอพยัคฆภูมิพิสัย จังหวัดมหาสารคาม ห่างจากตัวอำเภอประมาณ 20 กิโลเมตร ชายทุ่งกุลาร้องไห้","zh":"玛哈沙拉堪府帕亚卡普米披赛县孟涛乡孟涛村，距县城约20公里，位于古拉龙海平原边缘。","en":"In Ban Mueang Tao, Mueang Tao subdistrict, Phayakkhaphum Phisai, Maha Sarakham, about 20 km from the district town on the edge of Thung Kula Ronghai."}','K-001'),
('history','{"th":"เจ้าพ่อศรีนครเตาคือใคร","zh":"西那空涛神是谁？","en":"Who was Chao Pho Si Nakhon Tao?"}',
 '{"th":"ตามตำนานคือเซียงสี หัวหน้าชาวกุยที่อพยพมาจากอัตปือ ช่วยจับช้างเผือกให้กรุงศรีอยุธยา ได้รับบรรดาศักดิ์เป็นหลวงและพระศรีนครเตาท้าวเธอ ชาวบ้านนับถือเป็นบรรพบุรุษผู้สร้างบ้านเมืองเตา","zh":"据传说，他是从阿速坡迁来的奎族首领谢西，曾为阿瑜陀耶捕获白象，受封为“銮”和“帕·西那空涛陶特”，村民奉他为建立孟涛的祖先。","en":"By legend he was Xiang Si, a Kui leader from Attapeu who caught the white elephants for Ayutthaya and was made Luang, then Phra Si Nakhon Tao. Villagers honour him as the founder of Mueang Tao."}','K-008'),
('history','{"th":"ทำไมเรียกว่าเมืองเตา","zh":"为什么叫孟涛？","en":"Why is it called Mueang Tao?"}',
 '{"th":"\"เตา\" หมายถึงเตาถลุงเหล็ก ชุมชนมีหลักฐานการถลุงเหล็กสมัยก่อนประวัติศาสตร์ที่หนองเตาเหล็ก ทางทิศตะวันตกเฉียงใต้ของหมู่บ้าน","zh":"“涛”意为炼铁炉。村庄西南方的铁炉塘有史前炼铁的证据。","en":"Tao means an iron-smelting furnace; prehistoric iron working is evident at Nong Tao Lek, south-west of the village."}','K-002'),
('ritual','{"th":"งานประจำปีของเจ้าพ่อจัดเมื่อไร","zh":"神祠的年度庆典什么时候举行？","en":"When is the annual festival?"}',
 '{"th":"จัดในเดือน 4 ตามจันทรคติของทุกปี วันที่แน่นอนประกาศโดยชุมชน","zh":"每年农历四月举行，确切日期由社区公布。","en":"Every year in the fourth lunar month; the community announces the exact date."}','K-011'),
('ritual','{"th":"ต้องเตรียมอะไรไปไหว้เจ้าพ่อ","zh":"祭拜需要准备什么？","en":"What should I bring to pay respect?"}',
 '{"th":"แต่งกายสุภาพ เตรียมดอกไม้ ธูป เทียน และอธิษฐานด้วยความสงบ","zh":"穿着得体，准备鲜花、香和蜡烛，安静地祈祷。","en":"Dress modestly, bring flowers, incense and candles, and pray quietly."}','K-014');

-- ---------------------------------------------------------------- survey
INSERT INTO survey_questions(section, text, sort) VALUES
('satisfaction','{"th":"ความพึงพอใจโดยรวมต่อแพลตฟอร์ม","zh":"对平台的整体满意度","en":"Overall satisfaction with the platform"}',1),
('satisfaction','{"th":"เนื้อหาช่วยให้เข้าใจประวัติและความเชื่อของชุมชน","zh":"内容帮助我了解社区的历史与信仰","en":"Content helped me understand local history and beliefs"}',2),
('usability','{"th":"ใช้งานง่าย ค้นหาข้อมูลได้รวดเร็ว","zh":"易于使用，能快速找到信息","en":"Easy to use and quick to find information"}',1),
('usability','{"th":"แผนที่และการนำทางใช้ได้จริง","zh":"地图和导航实用","en":"Map and navigation were useful"}',2),
('ai_quality','{"th":"คำตอบของผู้ช่วย AI ถูกต้องและน่าเชื่อถือ","zh":"AI助手的回答准确可信","en":"The AI assistant answers were accurate and trustworthy"}',1),
('ai_quality','{"th":"ผู้ช่วย AI แสดงแหล่งอ้างอิงชัดเจน","zh":"AI助手清楚地显示了参考来源","en":"The AI assistant showed its sources clearly"}',2);

-- ---------------------------------------------------------------- AI test set (starter rows)
INSERT INTO eval_questions(code, lang, question, reference_answer, relevant_docs, category) VALUES
('T-001','th','เจ้าพ่อศรีนครเตาอยู่ที่ไหน','บ้านเมืองเตา อ.พยัคฆภูมิพิสัย จ.มหาสารคาม','{K-001}','history'),
('T-002','zh','祭拜需要准备什么？','鲜花、香、蜡烛','{K-014}','ritual'),
('T-003','th','พิธีขอฝนจัดเมื่อไหร่','ก่อนฤดูทำนา','{K-015}','ritual'),
('T-004','th','ทำไมชื่อเมืองเตา','เตาถลุงเหล็ก หนองเตาเหล็ก','{K-002}','history'),
('T-005','th','ดาบนวโลหะทำจากอะไร','โลหะเก้าชนิด','{K-007}','legend'),
('T-006','zh','西那空涛神是谁？','谢西，奎族首领，受封为銮·西那空涛','{K-008}','legend'),
('T-007','th','คำขวัญตำบลเมืองเตาคืออะไร','ถิ่นข้าวหอมมะลิหนึ่งศูนย์ห้า ... นับถือเจ้าพ่อศรีนครเตา','{K-012}','history'),
('T-008','en','What artefacts were found at Ban Mueang Tao?','Potsherds, beads, terracotta amulets, bangles, stone tools','{K-004}','history');

-- ---------------------------------------------------------------- images
-- Illustrative stock photos (Unsplash licence). Replace with the community own photos (survey item A1/A3).
UPDATE places SET images = ARRAY['https://images.unsplash.com/photo-1758346974564-07a164871e7d?auto=format&fit=crop&w=1200&q=75','https://images.unsplash.com/photo-1758346974892-be6bc5b36815?auto=format&fit=crop&w=1200&q=75'] WHERE code = 'P-001';
UPDATE places SET images = ARRAY['https://images.unsplash.com/photo-1562835156-75d18f5ef57f?auto=format&fit=crop&w=1200&q=75','https://images.unsplash.com/photo-1679429610125-c3e3278cd85f?auto=format&fit=crop&w=1200&q=75'] WHERE code = 'P-002';
UPDATE places SET images = ARRAY['https://images.unsplash.com/photo-1747396108528-682b02327818?auto=format&fit=crop&w=1200&q=75'] WHERE code = 'P-003';
UPDATE places SET images = ARRAY['https://images.unsplash.com/photo-1647167363159-1b3693fc468b?auto=format&fit=crop&w=1200&q=75'] WHERE code = 'P-004';
UPDATE knowledge SET cover_url = 'https://images.unsplash.com/photo-1679429610125-c3e3278cd85f?auto=format&fit=crop&w=1200&q=75' WHERE doc_code IN ('K-001','K-003','K-004','K-013');
UPDATE knowledge SET cover_url = 'https://images.unsplash.com/photo-1750433835485-b1be2eacc19a?auto=format&fit=crop&w=1200&q=75' WHERE doc_code IN ('K-005','K-006','K-009','K-011');
UPDATE knowledge SET cover_url = 'https://images.unsplash.com/photo-1699941325571-e424dada03a5?auto=format&fit=crop&w=1200&q=75' WHERE doc_code IN ('K-002','K-007','K-012','K-014');
UPDATE knowledge SET cover_url = 'https://images.unsplash.com/photo-1757059112718-badebf0d85ca?auto=format&fit=crop&w=1200&q=75' WHERE doc_code IN ('K-008','K-010','K-015');
UPDATE events SET cover_url = 'https://images.unsplash.com/photo-1757059112718-badebf0d85ca?auto=format&fit=crop&w=1200&q=75';

-- ---------------------------------------------------------------- news
INSERT INTO news(title, body, cover_url, is_banner, published_at) VALUES
('{"th":"มรภ.มหาสารคาม พัฒนาแพลตฟอร์มท่องเที่ยวเชิงศรัทธา \"เจ้าพ่อศรีนครเตา\"","zh":"玛哈沙拉堪皇家大学开发“西那空涛神”信仰旅游平台","en":"Rajabhat Maha Sarakham builds a faith tourism platform for Chao Pho Si Nakhon Tao"}',
 '{"th":"โครงการวิจัยภายใต้ทุนมูลฐาน (Fundamental Fund) ปีงบประมาณ 2569 ของมหาวิทยาลัยราชภัฏมหาสารคาม พัฒนาแพลตฟอร์มท่องเที่ยวเชิงศรัทธาที่รองรับภาษาไทยและภาษาจีน รวบรวมประวัติ ตำนาน และพิธีกรรมของเจ้าพ่อศรีนครเตา เพื่อส่งเสริมการท่องเที่ยวและเศรษฐกิจของชุมชนบ้านเมืองเตา","zh":"玛哈沙拉堪皇家大学2569财年基础研究基金项目开发了支持泰语和汉语的信仰旅游平台，汇集西那空涛神的历史、传说与仪式，以促进孟涛村的旅游和社区经济。","en":"A Fundamental Fund project (fiscal 2026) at Rajabhat Maha Sarakham University is building a Thai and Chinese faith tourism platform that gathers the history, legends and rites of Chao Pho Si Nakhon Tao to support tourism and the local economy of Ban Mueang Tao."}',
 'https://images.unsplash.com/photo-1789971649652-ae5dfeeaa7b4?auto=format&fit=crop&w=1200&q=75', true, now() - interval '20 days'),
('{"th":"คณะวิจัยเก็บข้อมูลประวัติชุมชนจากปราชญ์ชาวบ้าน","zh":"研究团队向社区学者收集村史资料","en":"Researchers record village history with community scholars"}',
 '{"th":"คณะผู้วิจัยลงพื้นที่บ้านเมืองเตาเพื่อสัมภาษณ์ผู้นำชุมชน ปราชญ์ชาวบ้าน และชาวบ้าน รวบรวมตำนานเจ้าพ่อศรีนครเตา ความเชื่อ และของดีของตำบล ข้อมูลทั้งหมดจะผ่านการตรวจสอบจากชุมชนก่อนเผยแพร่","zh":"研究团队走访孟涛村，采访社区领袖、民间学者和村民，收集西那空涛神的传说、信仰和本乡特色。所有资料在发布前将经社区核实。","en":"The research team visited Ban Mueang Tao to interview community leaders, local scholars and villagers about the legend, beliefs and local products. Everything is checked by the community before it is published."}',
 'https://images.unsplash.com/photo-1750433835485-b1be2eacc19a?auto=format&fit=crop&w=1200&q=75', false, now() - interval '45 days');

-- ---------------------------------------------------------------- local products (from the subdistrict motto)
-- Shops and community products are created by migration 003_community_shops.sql.
