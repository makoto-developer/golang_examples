create table if not exists online_shop.orders
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

comment on table online_shop.orders is '注文履歴';

comment on column online_shop.orders.order_item_group_id is '商品IDのリスト';

comment on column online_shop.orders.user_id is 'ユーザID';

comment on column online_shop.orders.amount is '税込価格(商品価格の合計)';

comment on column online_shop.orders.amount_without_tax is '税抜価格(税抜の商品価格の合計)';

comment on column online_shop.orders.tax is '消費税(商品価格に対する税の合計)';

comment on column online_shop.orders.tax_rate_id is '消費税率ID';

comment on column online_shop.orders.why_deleted is '削除した理由';

create index orders_user_id_index
    on online_shop.orders (user_id);

create index orders_user_id_orders_item_group_id_index
    on online_shop.orders (user_id, order_item_group_id);

create index orders_user_id_id_index
    on online_shop.orders (user_id, id);

create index orders_order_item_group_id_index
    on online_shop.orders (order_item_group_id);

