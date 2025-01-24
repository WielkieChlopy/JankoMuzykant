defmodule JankoMuzykant.Youtube do
  def get_song_play_url(url) do
    {output, status} = System.shell("yt-dlp -f bestaudio -g \"#{url}\"")
    if status == 0 do
      {:ok, output}
    else
      {:error, "Failed to load song URL"}
    end
  end
end
