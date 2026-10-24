
-- Illustrative stock photos (Unsplash licence). Replace with the community own photos (survey item A1/A3).
UPDATE places SET images = ARRAY['https://images.unsplash.com/photo-1758346974564-07a164871e7d?auto=format&fit=crop&w=1200&q=75','https://images.unsplash.com/photo-1758346974892-be6bc5b36815?auto=format&fit=crop&w=1200&q=75'] WHERE code = 'P-001';
UPDATE places SET images = ARRAY['https://images.unsplash.com/photo-1562835156-75d18f5ef57f?auto=format&fit=crop&w=1200&q=75','https://images.unsplash.com/photo-1679429610125-c3e3278cd85f?auto=format&fit=crop&w=1200&q=75'] WHERE code = 'P-002';
UPDATE places SET images = ARRAY['https://images.unsplash.com/photo-1747396108528-682b02327818?auto=format&fit=crop&w=1200&q=75'] WHERE code = 'P-003';
UPDATE places SET images = ARRAY['https://images.unsplash.com/photo-1647167363159-1b3693fc468b?auto=format&fit=crop&w=1200&q=75'] WHERE code = 'P-004';
UPDATE knowledge SET cover_url = 'https://images.unsplash.com/photo-1679429610125-c3e3278cd85f?auto=format&fit=crop&w=1200&q=75' WHERE doc_code = 'K-001';
UPDATE knowledge SET cover_url = 'https://images.unsplash.com/photo-1750433835485-b1be2eacc19a?auto=format&fit=crop&w=1200&q=75' WHERE doc_code = 'K-002';
UPDATE knowledge SET cover_url = 'https://images.unsplash.com/photo-1699941325571-e424dada03a5?auto=format&fit=crop&w=1200&q=75' WHERE doc_code = 'K-003';
UPDATE knowledge SET cover_url = 'https://images.unsplash.com/photo-1757059112718-badebf0d85ca?auto=format&fit=crop&w=1200&q=75' WHERE doc_code = 'K-004';
UPDATE events SET cover_url = 'https://images.unsplash.com/photo-1757059112718-badebf0d85ca?auto=format&fit=crop&w=1200&q=75';
INSERT INTO news(title, body, cover_url, is_banner, published_at) VALUES
('{"th":"เตรียมพบงานบุญประจำปีของชุมชน (ตัวอย่าง)","zh":"社区年度功德节即将到来（示例）","en":"The annual merit festival of the village is coming (sample)"}',
 '{"th":"[ตัวอย่าง] ขอเชิญร่วมงานบุญประจำปีและพิธีสักการะเจ้าพ่อศรีนครเตา รายละเอียดวันเวลาจะประกาศโดยชุมชน","zh":"【示例】诚邀参加年度功德节与西那空涛神祭拜仪式，具体日期由社区公布。","en":"[Sample] Join the annual merit festival and offering ceremony. Dates will be announced by the community."}',
 'https://images.unsplash.com/photo-1789971649652-ae5dfeeaa7b4?auto=format&fit=crop&w=1200&q=75', true, now());
INSERT INTO shops(name, kind, products, description, lat, lon, contact, images) VALUES
('{"th":"กลุ่มทอผ้าบ้านเมืองเตา (ตัวอย่าง)","zh":"孟涛村织布小组（示例）","en":"Ban Mueang Tao weaving group (sample)"}','product',
 '{"th":"ผ้าไหม ผ้าฝ้ายทอมือ","zh":"手工丝绸、棉布","en":"Hand-woven silk and cotton"}',
 '{"th":"[ตัวอย่าง] ผลิตภัณฑ์ทอมือของกลุ่มแม่บ้าน","zh":"【示例】妇女小组的手工织品","en":"[Sample] Hand-woven goods by the village women group"}',
 15.5150, 103.1925, '', ARRAY['https://images.unsplash.com/photo-1569909115134-a0426936c879?auto=format&fit=crop&w=1200&q=75']),
('{"th":"ครัวอีสานบ้านเมืองเตา (ตัวอย่าง)","zh":"孟涛村东北厨房（示例）","en":"Mueang Tao Isan kitchen (sample)"}','food',
 '{"th":"ส้มตำ ไก่ย่าง ข้าวเหนียว","zh":"青木瓜沙拉、烤鸡、糯米饭","en":"Som tam, grilled chicken, sticky rice"}',
 '{"th":"[ตัวอย่าง] อาหารอีสานพื้นบ้าน","zh":"【示例】地道东北菜","en":"[Sample] Home-style Isan food"}',
 15.5135, 103.1958, '', ARRAY['https://images.unsplash.com/photo-1647167363124-d3567bdcb7df?auto=format&fit=crop&w=1200&q=75']);
