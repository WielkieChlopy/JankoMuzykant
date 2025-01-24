defmodule JankoMuzykantWeb.SongsController do
  use JankoMuzykantWeb, :controller

  def songs(conn, _params) do
    render(conn, :songs)
  end


end
