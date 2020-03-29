ALTER TABLE  `red_envelope_goods` 
ADD COLUMN `origin_envelope_no` varchar(32) NULL DEFAULT '' COMMENT '原订单号' AFTER `updated_at`;