-- Shop details (27 Oct 2026): price range, unit, how to buy, shipping, the maker's story and
-- a list of sellers, so one product can be sold by many households (project goal: 20).
-- Sellers are public: admins enter only people who agreed to be listed.

ALTER TABLE shops
    ADD COLUMN price_from integer CHECK (price_from >= 0),
    ADD COLUMN price_to   integer CHECK (price_to >= 0),
    ADD COLUMN price_unit jsonb,
    ADD COLUMN how_to_buy jsonb,
    ADD COLUMN story      jsonb,
    ADD COLUMN ships      boolean NOT NULL DEFAULT false,
    ADD COLUMN sellers    jsonb   NOT NULL DEFAULT '[]' CHECK (jsonb_typeof(sellers) = 'array');

-- Illustrative stock photos (Unsplash licence) until the community supplies its own.
-- Only rows that still have no photo, or only the earlier placeholder photos, are changed.
UPDATE shops s
SET images = v.images, updated_at = now()
FROM (VALUES
('ผ้าไหมมัดหมี่เมืองเตา', ARRAY['https://images.unsplash.com/photo-1569909115134-a0426936c879?auto=format&fit=crop&w=1200&q=75','https://images.unsplash.com/photo-1719154717733-d37860cdee24?auto=format&fit=crop&w=1200&q=75']),
('หมวกจักสานเมืองเตา', ARRAY['https://images.unsplash.com/photo-1505471768190-275e2ad7b3f9?auto=format&fit=crop&w=1200&q=75']),
('เครื่องจักสานไม้ไผ่บ้านเมืองเตา', ARRAY['https://images.unsplash.com/photo-1771759539572-0b94498364b6?auto=format&fit=crop&w=1200&q=75']),
('ข้าวหอมมะลิ 105 ทุ่งกุลาร้องไห้', ARRAY['https://images.unsplash.com/photo-1536153635972-1fc2e818f642?auto=format&fit=crop&w=1200&q=75']),
('สมุนไพรพื้นบ้านเมืองเตา', ARRAY['https://images.unsplash.com/photo-1560705889-dc78ff46bbd6?auto=format&fit=crop&w=1200&q=75']),
('อาหารพื้นบ้านอีสานเมืองเตา', ARRAY['https://images.unsplash.com/photo-1506781961370-37a89d6b3095?auto=format&fit=crop&w=1200&q=75','https://images.unsplash.com/photo-1758346975340-b817918e04cd?auto=format&fit=crop&w=1200&q=75','https://images.unsplash.com/photo-1647167363124-d3567bdcb7df?auto=format&fit=crop&w=1200&q=75']),
('ชุดเครื่องสักการะเจ้าพ่อศรีนครเตา', ARRAY['https://images.unsplash.com/photo-1780247584867-d332ed89d423?auto=format&fit=crop&w=1200&q=75','https://images.unsplash.com/photo-1782309088131-e376ec2487de?auto=format&fit=crop&w=1200&q=75']),
('กิจกรรมเรียนรู้ ทอผ้าและจักสานกับปราชญ์ชุมชน', ARRAY['https://images.unsplash.com/photo-1719154717733-d37860cdee24?auto=format&fit=crop&w=1200&q=75','https://images.unsplash.com/photo-1771759539572-0b94498364b6?auto=format&fit=crop&w=1200&q=75'])
) AS v(name_th, images)
WHERE s.name->>'th' = v.name_th
  AND s.images <@ ARRAY['https://images.unsplash.com/photo-1569909115134-a0426936c879?auto=format&fit=crop&w=1200&q=75','https://images.unsplash.com/photo-1747396108528-682b02327818?auto=format&fit=crop&w=1200&q=75','https://images.unsplash.com/photo-1647167363124-d3567bdcb7df?auto=format&fit=crop&w=1200&q=75']::text[];
