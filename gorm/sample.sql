-- postgresql用のコードです。テーブルは思いつきで作ったやつなので本番利用しないで。
create database myshop
create schema online_shop;
create table "orders"
(
    id                  bigserial
        constraint orders_pk
            primary key,
    order_item_group_id bigserial,
    user_id             bigserial,
    amount              bigserial,
    amount_without_tax  bigserial,
    tax                 bigserial,
    tax_rate_id         smallserial,
    created_at          timestamp default CURRENT_TIMESTAMP,
    updated_at          timestamp default CURRENT_TIMESTAMP,
    deleted_at          timestamp,
    why_deleted         text
);

comment on table "orders" is '注文履歴';

comment on column "orders".order_item_group_id is '商品IDのリスト';

comment on column "orders".user_id is 'ユーザID';

comment on column "orders".amount is '税込価格(商品価格の合計)';

comment on column "orders".amount_without_tax is '税抜価格(税抜の商品価格の合計)';

comment on column "orders".tax is '消費税(商品価格に対する税の合計)';

comment on column "orders".tax_rate_id is '消費税率ID';

comment on column "orders".why_deleted is '削除した理由';

create index orders_user_id_index
    on "orders" (user_id);

create index orders_user_id_orders_item_group_id_index
    on "orders" (user_id, orders_item_group_id);

create index orders_user_id_id_index
    on "orders" (user_id, id);

create index orders_order_item_group_id_index
    on "orders" (order_item_group_id);

