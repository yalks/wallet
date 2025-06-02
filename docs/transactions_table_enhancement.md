# Transactions 表增强建议

## 概述

当前的 `transactions` 表缺少记录用户发起交易时的完整请求信息。为了更好地追踪用户行为、进行审计和问题排查，建议添加以下字段来记录用户请求的完整数据。

## 新增字段建议

### 1. 用户请求原始信息

```sql
-- 用户请求的原始金额 (用户输入的金额)
`request_amount` decimal(36,8) DEFAULT NULL COMMENT '用户请求的原始金额 (用户输入的金额)',

-- 用户请求的参考信息 (如转账备注、提现地址等)
`request_reference` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户请求的参考信息 (如转账备注、提现地址等)',

-- 用户请求的元数据 (JSON格式存储扩展信息)
`request_metadata` json DEFAULT NULL COMMENT '用户请求的元数据 (JSON格式存储扩展信息)',
```

### 2. 请求来源和环境信息

```sql
-- 请求来源 (telegram, web, api, admin等)
`request_source` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '请求来源 (telegram, web, api, admin等)',

-- 用户请求的IP地址
`request_ip` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户请求的IP地址',

-- 用户请求的User-Agent
`request_user_agent` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户请求的User-Agent',
```

### 3. 时间追踪

```sql
-- 用户发起请求的时间戳
`request_timestamp` timestamp NULL DEFAULT NULL COMMENT '用户发起请求的时间戳',

-- 交易处理完成时间
`processed_at` timestamp NULL DEFAULT NULL COMMENT '交易处理完成时间',
```

### 4. 费用和汇率信息

```sql
-- 手续费金额
`fee_amount` decimal(36,8) DEFAULT '0.00000000' COMMENT '手续费金额',

-- 手续费类型 (fixed, percentage)
`fee_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '手续费类型 (fixed, percentage)',

-- 汇率 (如果涉及币种转换)
`exchange_rate` decimal(36,8) DEFAULT NULL COMMENT '汇率 (如果涉及币种转换)',
```

### 5. 目标用户信息 (用于转账、红包等)

```sql
-- 目标用户ID (转账、红包等操作的接收方)
`target_user_id` int unsigned DEFAULT NULL COMMENT '目标用户ID (转账、红包等操作的接收方)',

-- 目标用户名 (转账、红包等操作的接收方用户名)
`target_username` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '目标用户名 (转账、红包等操作的接收方用户名)',
```

## 完整的 ALTER TABLE 语句

```sql
ALTER TABLE `transactions` 
ADD COLUMN `request_amount` decimal(36,8) DEFAULT NULL COMMENT '用户请求的原始金额 (用户输入的金额)' AFTER `business_id`,
ADD COLUMN `request_reference` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户请求的参考信息 (如转账备注、提现地址等)' AFTER `request_amount`,
ADD COLUMN `request_metadata` json DEFAULT NULL COMMENT '用户请求的元数据 (JSON格式存储扩展信息)' AFTER `request_reference`,
ADD COLUMN `request_source` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '请求来源 (telegram, web, api, admin等)' AFTER `request_metadata`,
ADD COLUMN `request_ip` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户请求的IP地址' AFTER `request_source`,
ADD COLUMN `request_user_agent` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户请求的User-Agent' AFTER `request_ip`,
ADD COLUMN `request_timestamp` timestamp NULL DEFAULT NULL COMMENT '用户发起请求的时间戳' AFTER `request_user_agent`,
ADD COLUMN `processed_at` timestamp NULL DEFAULT NULL COMMENT '交易处理完成时间' AFTER `request_timestamp`,
ADD COLUMN `fee_amount` decimal(36,8) DEFAULT '0.00000000' COMMENT '手续费金额' AFTER `processed_at`,
ADD COLUMN `fee_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '手续费类型 (fixed, percentage)' AFTER `fee_amount`,
ADD COLUMN `exchange_rate` decimal(36,8) DEFAULT NULL COMMENT '汇率 (如果涉及币种转换)' AFTER `fee_type`,
ADD COLUMN `target_user_id` int unsigned DEFAULT NULL COMMENT '目标用户ID (转账、红包等操作的接收方)' AFTER `exchange_rate`,
ADD COLUMN `target_username` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '目标用户名 (转账、红包等操作的接收方用户名)' AFTER `target_user_id`;
```

## 建议的索引

```sql
-- 请求来源索引
ALTER TABLE `transactions` ADD INDEX `idx_request_source` (`request_source`);

-- 请求时间索引
ALTER TABLE `transactions` ADD INDEX `idx_request_timestamp` (`request_timestamp`);

-- 目标用户索引
ALTER TABLE `transactions` ADD INDEX `idx_target_user` (`target_user_id`);

-- 费用类型索引
ALTER TABLE `transactions` ADD INDEX `idx_fee_type` (`fee_type`);

-- 复合索引：用户+请求来源+时间
ALTER TABLE `transactions` ADD INDEX `idx_user_source_time` (`user_id`, `request_source`, `request_timestamp`);
```

## 字段说明和使用场景

### request_amount vs amount
- `request_amount`: 用户输入的原始金额
- `amount`: 实际处理的金额（可能包含手续费计算后的结果）

### request_metadata 存储示例
```json
{
  "red_packet_id": 123,
  "recipient": "username",
  "transfer_id": 456,
  "withdraw_address": "TXXXxxxXXXxxxXXX",
  "device_info": {
    "platform": "telegram",
    "version": "1.0.0"
  }
}
```

### request_source 可能的值
- `telegram`: Telegram Bot
- `web`: Web界面
- `api`: API调用
- `admin`: 管理员操作
- `system`: 系统自动操作

## 数据迁移注意事项

1. **现有数据处理**: 对于现有记录，新字段将为 NULL
2. **应用程序更新**: 需要同步更新相关的 Go 结构体和业务逻辑
3. **性能影响**: 新增字段较多，建议分批执行 ALTER TABLE
4. **存储空间**: JSON 字段和 TEXT 字段会增加存储需求

## 业务价值

1. **完整的审计追踪**: 记录用户请求的完整信息
2. **问题排查**: 能够追溯到用户的原始请求
3. **数据分析**: 分析用户行为模式和请求来源
4. **合规要求**: 满足金融监管的数据记录要求
5. **用户体验**: 更好地理解用户操作流程
