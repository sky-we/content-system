CREATE DATABASE `cms_content` DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

use cms_content;
CREATE TABLE `content_details` (
                                     `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增ID',
                                     `title` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '内容标题',
                                     `description` text COLLATE utf8mb4_unicode_ci COMMENT '内容描述',
                                     `author` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '作者',
                                     `video_url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '视频播放URL',
                                     `thumbnail` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '封面图URL',
                                     `category` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '内容分类',
                                     `duration` bigint DEFAULT '0' COMMENT '内容时长',
                                     `resolution` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '分辨率 如720p、1080p',
                                     `file_size` bigint DEFAULT '0' COMMENT '文件大小',
                                     `format` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '文件格式 如MP4、AVI',
                                     `quality` int DEFAULT '0' COMMENT '视频质量 1-高清 2-标清',
                                     `approval_status` int DEFAULT '1' COMMENT '审核状态 1-审核中 2-审核通过 3-审核不通过',
                                     `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '内容更新时间',
                                     `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '内容创建时间',
                                     PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='内容详情';



CREATE TABLE IF NOT EXISTS t_idx_content_details (
                                                     id BIGINT NOT NULL AUTO_INCREMENT,
                                                     content_id VARCHAR(255) DEFAULT '' COMMENT '内容ID',
    title VARCHAR(255) DEFAULT '' COMMENT '内容标题',
    author VARCHAR(255) DEFAULT '' COMMENT '作者',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '内容更新时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '内容创建时间',
    PRIMARY KEY (id)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;