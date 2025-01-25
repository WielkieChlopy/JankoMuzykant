defmodule JankoMuzykantWeb.PlayerLive do
  use JankoMuzykantWeb, :live_view

  import JankoMuzykant.Youtube

  def mount(_params, _session, socket) do
    {:ok,
     assign(socket,
       songs: [
         %{id: 1, loading: false, duration: nil, title: "Song 1", url: "https://www.youtube.com/watch?v=dQw4w9WgXcQ"},
         %{id: 2, loading: false, duration: nil, title: "Song 2", url: "https://www.youtube.com/watch?v=qwVp49uwHgk"},
       ],
       current_song: nil,
       next_song_id: 1,
       song_play_url: nil,
       next_song_url: nil
     )}
  end

  def handle_event("play", %{"id" => id}, socket) do
    song = Enum.find(socket.assigns.songs, fn song -> song.id == String.to_integer(id) end)

    if {:ok, output} = get_song_play_url(song.url) do
      {:noreply, assign(socket, current_song: song, song_play_url: output, next_song_url: nil)}
    else
      {:noreply, put_flash(socket, :error, "Could not load song URL")}
    end
  end

  def handle_event("play_next_song", _params, socket) do
    IO.puts("play_next_song")
    {:noreply, assign(socket, current_song: socket.assigns.next_song, song_play_url: socket.assigns.next_song_url, next_song_url: nil)}
  end

  def handle_event("load_next_song", _params, socket) do
    IO.puts("load_next_song")
    current_index =
      Enum.find_index(socket.assigns.songs, fn s -> s.id == socket.assigns.current_song.id end)

    next_song = Enum.at(socket.assigns.songs, current_index + 1)

    if next_song do
      if {:ok, output} = get_song_play_url(next_song.url) do
        {:noreply, assign(socket, next_song_url: output, next_song: next_song)}
      else
        {:noreply, socket}
      end
    else
      {:noreply, socket}
    end
  end

  def handle_event("add_song", %{"title" => title, "url" => url}, socket) do
    {:noreply, assign(socket, songs: [%{id: socket.assigns.next_song_id, title: title, url: url} | socket.assigns.songs], next_song_id: socket.assigns.next_song_id + 1)}
  end



  def render(assigns) do
    ~H"""
    <div class="max-w-lg mx-auto mt-8">
      <h1 class="text-2xl font-bold mb-4">Audio Player</h1>
      <div class="bg-gray-100 p-4 rounded flex justify-between items-center">
        <form phx-submit="add_song">
          <input type="text" name="title" placeholder="Title" />
          <input type="text" name="url" placeholder="URL" />
          <button type="submit">Add Song</button>
        </form>
      </div>
      <div class="bg-gray-100 p-4 rounded">
        <h2 class="font-semibold mb-2">Queue</h2>
        <ul>
          <%= for song <- @songs do %>
            <li class="flex justify-between items-center py-2">
              <span><%= song.title %></span>
              <button
                phx-click="play"
                phx-value-id={song.id}
                class="bg-blue-500 text-white px-3 py-1 rounded"
              >
                Play
              </button>
            </li>
          <% end %>
        </ul>
      </div>

      <div>
        <%= if @current_song do %>
          <%= @current_song.title %>
          <div id="audio-player-container">
            <audio phx-hook="AudioPlayer" id="audioPlayer" controls autoplay>
              <source src={@song_play_url} type="audio/mp3" />
              Your browser does not support the audio element.
            </audio> 
          </div>
        <% end %>
      </div>
    </div>

    """
  end
end
