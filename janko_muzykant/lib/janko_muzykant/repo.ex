defmodule JankoMuzykant.Repo do
  use Ecto.Repo,
    otp_app: :janko_muzykant,
    adapter: Ecto.Adapters.Postgres
end
