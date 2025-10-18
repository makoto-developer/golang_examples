-- postgresql用のコードです。テーブルは思いつきで作ったやつなので本番利用しないで。
create database myshop
create schema online_shop;
create table "order"
(
    id                  bigserial
        constraint order_pk
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

comment on table "order" is '注文履歴';

comment on column "order".order_item_group_id is '商品IDのリスト';

comment on column "order".user_id is 'ユーザID';

comment on column "order".amount is '税込価格(商品価格の合計)';

comment on column "order".amount_without_tax is '税抜価格(税抜の商品価格の合計)';

comment on column "order".tax is '消費税(商品価格に対する税の合計)';

comment on column "order".tax_rate_id is '消費税率ID';

comment on column "order".why_deleted is '削除した理由';

alter table "order"
    owner to "FJIO8880awz0";

create index order_user_id_index
    on "order" (user_id);

create index order_order_item_group_id_index
    on "order" (order_item_group_id);

create index order_user_id_order_item_group_id_index
    on "order" (user_id, order_item_group_id);

create index order_user_id_id_index
    on "order" (user_id, id);

