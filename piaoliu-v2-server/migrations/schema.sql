-- 漂流瓶数据库初始化脚本
-- MySQL 8.0+
-- 注意：必须 SET NAMES utf8mb4，否则中文备注在 docker 初始化时会被按 latin1 二次编码变成乱码

CREATE DATABASE IF NOT EXISTS drift_bottle DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
USE drift_bottle;
SET NAMES utf8mb4;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
  id BIGINT NOT NULL AUTO_INCREMENT COMMENT '用户主键ID',
  openid VARCHAR(64) NOT NULL COMMENT '微信 openid，用户唯一标识',
  unionid VARCHAR(64) DEFAULT NULL COMMENT '微信 unionid（同主体跨应用统一ID）',
  nickname VARCHAR(50) NOT NULL COMMENT '用户昵称',
  avatar_url VARCHAR(255) DEFAULT NULL COMMENT '头像 URL',
  province_code VARCHAR(20) NOT NULL COMMENT '省份编码',
  province_name VARCHAR(50) NOT NULL COMMENT '省份名称',
  city_code VARCHAR(20) NOT NULL COMMENT '城市编码',
  city_name VARCHAR(50) NOT NULL COMMENT '城市名称',
  district_code VARCHAR(20) DEFAULT NULL COMMENT '区县编码',
  district_name VARCHAR(50) DEFAULT NULL COMMENT '区县名称',
  status TINYINT DEFAULT 1 COMMENT '账号状态：1正常 2禁言 3封号',
  address_updated_at DATETIME(3) DEFAULT NULL COMMENT '最近一次修改地址时间',
  created_at DATETIME(3) DEFAULT NULL COMMENT '创建时间',
  updated_at DATETIME(3) DEFAULT NULL COMMENT '更新时间',
  gender TINYINT DEFAULT 0 COMMENT '性别：0未知 1男 2女',
  birthday VARCHAR(20) DEFAULT NULL COMMENT '生日（YYYY-MM-DD）',
  bottle_count BIGINT DEFAULT 0 COMMENT '累计扔出的瓶子数',
  interest_tags JSON DEFAULT NULL COMMENT '兴趣标签数组',
  PRIMARY KEY (id),
  UNIQUE KEY uk_users_open_id (openid),
  KEY idx_openid (openid),
  KEY idx_city_code (city_code),
  KEY idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 瓶子表
CREATE TABLE IF NOT EXISTS bottles (
  id BIGINT NOT NULL AUTO_INCREMENT COMMENT '瓶子主键ID',
  user_id BIGINT NOT NULL COMMENT '发布者用户ID',
  content TEXT COMMENT '瓶内文字内容',
  media_urls JSON DEFAULT NULL COMMENT '图片/语音 URL 数组',
  bottle_type TINYINT DEFAULT 1 COMMENT '瓶子类型：1普通瓶 2语音瓶 3秘密瓶',
  province_code VARCHAR(20) NOT NULL COMMENT '发布地省份编码',
  city_code VARCHAR(20) NOT NULL COMMENT '发布地城市编码',
  district_code VARCHAR(20) DEFAULT NULL COMMENT '发布地区县编码',
  status TINYINT DEFAULT 0 COMMENT '审核状态：0待审核 1已通过 2未通过 3已删除',
  is_picked TINYINT(1) DEFAULT 0 COMMENT '是否已被捡过',
  picked_count BIGINT DEFAULT 0 COMMENT '被捡次数',
  visible_scope TINYINT DEFAULT 1 COMMENT '可见范围：1同城 2同省',
  expire_at DATETIME(3) DEFAULT NULL COMMENT '秘密瓶过期时间',
  created_at DATETIME(3) DEFAULT NULL COMMENT '创建时间',
  updated_at DATETIME(3) DEFAULT NULL COMMENT '更新时间',
  tags JSON DEFAULT NULL COMMENT '瓶子标签数组',
  PRIMARY KEY (id),
  KEY idx_user_id (user_id),
  KEY idx_city_status (city_code, status),
  KEY idx_district_status (district_code, status),
  KEY idx_created_at (created_at),
  KEY idx_expire_at (expire_at),
  CONSTRAINT fk_bottles_user FOREIGN KEY (user_id) REFERENCES users (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='瓶子表';

-- 捡瓶子记录表
CREATE TABLE IF NOT EXISTS bottle_picked_records (
  id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  bottle_id BIGINT NOT NULL COMMENT '被捡的瓶子ID',
  user_id BIGINT NOT NULL COMMENT '捡瓶人用户ID',
  is_replied TINYINT(1) DEFAULT 0 COMMENT '捡瓶人是否已回复',
  created_at DATETIME(3) DEFAULT NULL COMMENT '捡瓶时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_bottle_user (bottle_id, user_id),
  KEY idx_user_id (user_id),
  KEY idx_bottle_id (bottle_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='捡瓶子记录表';

-- 瓶子回复表
CREATE TABLE IF NOT EXISTS bottle_replies (
  id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  bottle_id BIGINT NOT NULL COMMENT '所属瓶子ID',
  user_id BIGINT NOT NULL COMMENT '回复人用户ID',
  content TEXT COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '回复内容',
  created_at DATETIME(3) DEFAULT NULL COMMENT '回复时间',
  PRIMARY KEY (id),
  KEY idx_bottle_replies_bottle_id (bottle_id),
  KEY idx_bottle_replies_user_id (user_id),
  CONSTRAINT fk_bottle_replies_user FOREIGN KEY (user_id) REFERENCES users (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='瓶子回复表';

-- 会话表
CREATE TABLE IF NOT EXISTS conversations (
  id BIGINT NOT NULL AUTO_INCREMENT COMMENT '会话主键ID',
  bottle_id BIGINT NOT NULL COMMENT '关联瓶子ID',
  user_a_id BIGINT NOT NULL COMMENT '扔瓶人用户ID',
  user_b_id BIGINT NOT NULL COMMENT '捡瓶/回复人用户ID',
  last_message_id BIGINT DEFAULT NULL COMMENT '最近一条消息ID',
  last_message_at DATETIME(3) DEFAULT NULL COMMENT '最近消息时间',
  is_closed TINYINT(1) DEFAULT 0 COMMENT '会话是否关闭',
  deleted_by_a TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'userA 侧是否已删除（仅对自己隐藏，新消息恢复）',
  deleted_by_b TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'userB 侧是否已删除（仅对自己隐藏，新消息恢复）',
  created_at DATETIME(3) DEFAULT NULL COMMENT '创建时间',
  updated_at DATETIME(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_bottle_pair (bottle_id, user_a_id, user_b_id),
  KEY idx_user_a (user_a_id),
  KEY idx_user_b (user_b_id),
  KEY idx_last_message_at (last_message_at),
  CONSTRAINT fk_conversations_bottle FOREIGN KEY (bottle_id) REFERENCES bottles (id),
  CONSTRAINT fk_conversations_last_msg FOREIGN KEY (last_message_id) REFERENCES messages (id),
  CONSTRAINT fk_conversations_user_a FOREIGN KEY (user_a_id) REFERENCES users (id),
  CONSTRAINT fk_conversations_user_b FOREIGN KEY (user_b_id) REFERENCES users (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会话表';

-- 消息表
CREATE TABLE IF NOT EXISTS messages (
  id BIGINT NOT NULL AUTO_INCREMENT COMMENT '消息主键ID',
  conversation_id BIGINT NOT NULL COMMENT '所属会话ID',
  sender_id BIGINT NOT NULL COMMENT '发送者用户ID',
  content TEXT COMMENT '消息文字内容',
  media_url VARCHAR(255) DEFAULT NULL COMMENT '图片/语音 URL',
  msg_type TINYINT DEFAULT 1 COMMENT '消息类型：1文字 2图片 3语音',
  is_read TINYINT(1) DEFAULT 0 COMMENT '是否已读',
  created_at DATETIME(3) DEFAULT NULL COMMENT '发送时间',
  PRIMARY KEY (id),
  KEY idx_conversation_id (conversation_id),
  KEY idx_sender_id (sender_id),
  KEY idx_created_at (created_at),
  CONSTRAINT fk_messages_sender FOREIGN KEY (sender_id) REFERENCES users (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='消息表';

-- 举报表
CREATE TABLE IF NOT EXISTS reports (
  id BIGINT NOT NULL AUTO_INCREMENT COMMENT '举报主键ID',
  reporter_id BIGINT NOT NULL COMMENT '举报人用户ID',
  target_user_id BIGINT NOT NULL COMMENT '被举报用户ID',
  target_bottle_id BIGINT DEFAULT NULL COMMENT '被举报瓶子ID',
  target_message_id BIGINT DEFAULT NULL COMMENT '被举报消息ID',
  reason TINYINT NOT NULL COMMENT '举报原因：1色情 2广告诈骗 3人身攻击 4违法违规 5其他',
  description VARCHAR(500) DEFAULT NULL COMMENT '举报补充说明',
  status TINYINT DEFAULT 0 COMMENT '处理状态：0待处理 1已处理-通过 2已处理-删除 3已处理-封号',
  handled_by BIGINT DEFAULT NULL COMMENT '处理人用户ID',
  handled_at DATETIME(3) DEFAULT NULL COMMENT '处理时间',
  created_at DATETIME(3) DEFAULT NULL COMMENT '举报时间',
  PRIMARY KEY (id),
  KEY idx_reporter (reporter_id),
  KEY idx_target_user (target_user_id),
  KEY idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='举报表';

-- 黑名单表
CREATE TABLE IF NOT EXISTS blacklists (
  id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  user_id BIGINT NOT NULL COMMENT '拉黑者用户ID',
  blocked_user_id BIGINT NOT NULL COMMENT '被拉黑用户ID',
  created_at DATETIME(3) DEFAULT NULL COMMENT '拉黑时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_user_blocked (user_id, blocked_user_id),
  KEY idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='黑名单表';

-- 签到记录表
CREATE TABLE IF NOT EXISTS checkin_records (
  id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  user_id BIGINT NOT NULL COMMENT '签到用户ID',
  consecutive_days BIGINT DEFAULT 1 COMMENT '连续签到天数',
  bottle_reward BIGINT DEFAULT 1 COMMENT '本次签到奖励瓶子数',
  created_at DATETIME(3) DEFAULT NULL COMMENT '签到时间',
  PRIMARY KEY (id),
  KEY idx_user_id (user_id),
  KEY idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='签到记录表';
