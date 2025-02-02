<script lang="ts">
	import * as Card from "@/components/ui/card/index.js";
	
	let { data } = $props()
	let playingSong = $state({
		url: '',
		play_url: '',
		title: '',
		duration: ''
	})
	let nextSong = $state({
		url: '',
		play_url: '',
		title: '',
		duration: ''
	})
	let tenSecondNotificationSent = $state(false)

	let audio: HTMLAudioElement | null = null

	async function loadNextSong() {
		const response = await fetch(`${data.backend_url}/api/v1/queue/next`, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json',
				'Authorization': `Bearer ${data.token}`
			},
		})
		nextSong = await response.json()
		console.log("load_next_song", nextSong)
	}

	async function playNextSong() {
		playingSong = nextSong
		if (audio) {
			audio.load()
		}
		console.log("play_next_song")
		const response = await fetch(`${data.backend_url}/api/v1/queue/next`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				'Authorization': `Bearer ${data.token}`
			},
		})
		console.log("play_next_song_response", response)
		//TODO: handle error
	}

	function onTimeUpdate(e: Event) {
		e.preventDefault()
		if (!audio) return
		const timeLeft = audio.duration - audio.currentTime
		
		if (timeLeft <= 10 && !tenSecondNotificationSent) {
			tenSecondNotificationSent = true
			console.log("load_next_song", { timeLeft })
			loadNextSong()
		}
	}


	function onEnded() {
		console.log("ended")
		playNextSong()
	}

	async function play_song(e: SubmitEvent) {
		e.preventDefault()
		const formData = new FormData(e.target as HTMLFormElement)
		const response = await fetch(`${data.backend_url}/api/v1/songs/details`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				'Authorization': `Bearer ${data.token}`
			},
			body: JSON.stringify({ url: formData.get('url') }),
		})
		playingSong = await response.json()
		tenSecondNotificationSent = false
		if (audio) {
			audio.load()
		}
	}
</script>

<div>
	<Card.Root>
		<Card.Header>
		  <Card.Title>Play Song now</Card.Title>
		</Card.Header>
		<Card.Content>
			<form onsubmit={play_song}>
				<input type="text" name="url" />
				<button type="submit">Add</button>
			</form>
		</Card.Content>
	  </Card.Root>
	  <Card.Root>
		<Card.Header>
		  <Card.Title>Player</Card.Title>
		  <Card.Description>Playing</Card.Description>
		</Card.Header>
		<Card.Content>
			<audio  bind:this={audio} src={playingSong.play_url} controls autoplay onended={onEnded} ontimeupdate={onTimeUpdate}></audio>
		</Card.Content>
	  </Card.Root>
	  <Card.Root>
		<Card.Header>
		  <Card.Title>Add to Queue</Card.Title>
		</Card.Header>
		<Card.Content>
			<form method="POST" action="?/add_song">
				<input type="text" name="url" />
				<button type="submit">Add</button>
			</form>
		</Card.Content>
	  </Card.Root>

	<Card.Root>
		<Card.Header>
		  <Card.Title>Queue</Card.Title>
		</Card.Header>
		<Card.Content>
			<ul>
				{#each data.queue as song}
					<li>{song.title}</li>
				{/each}
			</ul>
		</Card.Content>
	  </Card.Root>

</div>
