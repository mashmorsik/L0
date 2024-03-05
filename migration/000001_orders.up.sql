create table if not exists public.order
(
    uid                text primary key,
    track_number       text                     not null,
    entry              text                     not null,
    locale             text                     not null check (length(locale) <= 3),
    internal_signature text,
    customer_id        text                     not null check (length(customer_id) <= 50),
    delivery_service   text                     not null check (length(delivery_service) <= 50),
    shard_key          text                     not null check (length(shard_key) <= 10),
    sm_id              bigint                   not null,
    date_created       timestamp with time zone not null,
    oof_shard          text                     not null check (length(oof_shard) <= 10)
);

create table if not exists public.delivery
(
    id       bigint generated always as identity primary key,
    order_id text not null references public.order (uid),
    name     text not null check (length(name) <= 50),
    phone    text not null check (length(phone) <= 12),
    zip      text not null check (length(zip) <= 15),
    city     text not null check (length(city) <= 50),
    address  text not null check (length(address) <= 150),
    region   text not null check (length(phone) <= 100),
    email    text not null check (length(email) <= 200)
);

create table if not exists public.payment
(
    id            bigint generated always as identity primary key,
    transaction   text        not null references public.order (uid),
    request_id    text        not null,
    currency      text        not null check (length(currency) <= 10),
    provider      text        not null check (length(provider) <= 50),
    amount        int         not null,
    payment_dt    timestamptz not null,
    bank          text        not null check (length(bank) <= 50),
    delivery_cost int         not null,
    goods_total   int         not null,
    custom_fee    int         not null
);

create table if not exists public.order_item
(
    id           bigint generated always as identity primary key,
    order_id     text    not null references public.order (uid),
    chrt_id      bigint  not null,
    track_number text    not null check (length(track_number) <= 50),
    price        numeric not null,
    rid          text    not null,
    name         text    not null check (length(name) <= 50),
    sale         int     not null,
    size         text    not null check (length(size) <= 10),
    count        int     not null,
    total_price  numeric not null,
    nm_id        int     not null,
    brand        text    not null check (length(brand) <= 50),
    status       int     not null
);