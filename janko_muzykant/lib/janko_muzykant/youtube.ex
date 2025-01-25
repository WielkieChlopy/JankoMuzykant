defmodule JankoMuzykant.Youtube do
  def get_song_play_url(url) do
    {output, status} = System.shell("yt-dlp -f bestaudio -g \"#{url}\"")
    if status == 0 do
      {:ok, output}
    else
      {:error, "Failed to load song URL"}
    end
  end

  def get_song_info(url) do
    {output, status} = System.shell("yt-dlp --get-title  --get-duration \"#{url}\"")
    if status == 0 do
      [title, duration] = String.split(output, "\n")
      {:ok, %{title: title, duration: duration}}
    else
      {:error, "Failed to load song info"}
    end
  end
end
