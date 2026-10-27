-- Shorter shop descriptions (27 Oct 2026). The shops page now shows a
-- "seller details coming soon" badge for rows without contact details, so the
-- "(being collected...)" sentence is dropped from the text. Only rows that admins
-- have not given a contact yet are touched.

UPDATE shops s
SET description = v.description::jsonb, updated_at = now()
FROM (VALUES
('ผ้าไหมมัดหมี่เมืองเตา','{"th": "“งามวิไลไหมมัดหมี่” ของดีตามคำขวัญตำบล ทอมือโดยชาวบ้านเมืองเตา ซื้อเป็นของฝาก หรือสวมใส่มาไหว้เจ้าพ่อ", "zh": "孟涛乡标语中的“扎染丝绸秀丽”，由村民手工织成，可作伴手礼，也可穿着前来参拜。", "en": "The “lovely mudmee silk” of the subdistrict motto, hand-woven by villagers. Take it home, or wear it to the shrine."}'),
('หมวกจักสานเมืองเตา','{"th": "“รู้กันดีหมวกจักสาน” งานฝีมือขึ้นชื่อของตำบล กันแดดได้จริง เบา พกกลับเป็นของฝากง่าย", "zh": "孟涛乡标语中的“编织帽远近闻名”，遮阳实用，轻巧便携。", "en": "The well-known woven hats of the subdistrict motto: real sun shade, light and easy to carry home."}'),
('เครื่องจักสานไม้ไผ่บ้านเมืองเตา','{"th": "ภูมิปัญญาจักสานไม้ไผ่ของชาวนาทุ่งกุลา เหลือผู้สืบทอดไม่กี่คน รับทำตามสั่ง", "zh": "图拉平原稻农的竹编技艺，传承人已不多，按订单制作。", "en": "Bamboo craft of the Thung Kula rice farmers. Few makers remain; pieces are made to order."}'),
('ข้าวหอมมะลิ 105 ทุ่งกุลาร้องไห้','{"th": "ข้าวจากนาทุ่งกุลาร้องไห้ อาชีพหลักของชาวบ้าน ข้าวหอมมะลิทุ่งกุลาร้องไห้ขึ้นทะเบียนสิ่งบ่งชี้ทางภูมิศาสตร์ (GI)", "zh": "产自古拉龙海平原稻田，是村民的主要生计。古拉龙海茉莉香米已注册地理标志（GI）。", "en": "Grown in the Thung Kula Rong Hai fields, the village''s main livelihood. Thung Kula Rong Hai Hom Mali is a registered GI."}'),
('สมุนไพรพื้นบ้านเมืองเตา','{"th": "“สมุนไพรเลื่องลือ” ภูมิปัญญาดูแลสุขภาพที่ชาวบ้านใช้คู่กับการแพทย์สมัยใหม่", "zh": "孟涛乡标语中的“草药声名远扬”，村民与现代医学并用的养生智慧。", "en": "The famed herbs of the subdistrict motto, used by villagers alongside modern medicine."}'),
('อาหารพื้นบ้านอีสานเมืองเตา','{"th": "อาหารบ้านๆ แบบทุ่งกุลา จากวัตถุดิบในชุมชน แวะชิมได้หลังสักการะ", "zh": "图拉平原的家常菜，选用社区食材，参拜后可顺道品尝。", "en": "Everyday Thung Kula cooking from village produce. Stop for a meal after the shrine."}'),
('ชุดเครื่องสักการะเจ้าพ่อศรีนครเตา','{"th": "ดอกไม้ ธูป เทียน จัดชุดโดยครัวเรือนในชุมชน ขายหน้าศาล รายได้อยู่ในชุมชน", "zh": "由社区家庭配好的鲜花香烛，在神祠前出售，收入留在社区。", "en": "Flowers, incense and candles made up by village households and sold at the shrine."}'),
('กิจกรรมเรียนรู้ ทอผ้าและจักสานกับปราชญ์ชุมชน','{"th": "ลองทอผ้า สานไม้ไผ่ และทำอาหารพื้นบ้านกับคนในชุมชน จองล่วงหน้า", "zh": "跟村民体验织布、竹编和家常菜，需提前预约。", "en": "Try weaving, bamboo craft and local cooking with villagers. Book ahead."}'),
('โฮมสเตย์บ้านเมืองเตา','{"th": "พักในบ้านชาวบ้าน สำหรับผู้มาร่วมงานเดือน 4 หรือพิธีขอฝน", "zh": "住在村民家中，适合参加四月庆典或祈雨仪式的游客。", "en": "Stay in a village home for the fourth-month festival or the rain ceremony."}')
) AS v(name_th, description)
WHERE s.name->>'th' = v.name_th AND s.contact = '';
