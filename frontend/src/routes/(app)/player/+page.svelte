<script lang="ts">
	import * as Card from "@/components/ui/card/index.js";
	
	let { data } = $props()
	let playingSong = $state({
		url: '',
		play_url: '',
		title: '',
		duration: ''
	})

	let audio: HTMLAudioElement | null = null

	async function play_song(e: SubmitEvent) {
		e.preventDefault()
		const formData = new FormData(e.target as HTMLFormElement)
		const response = await fetch(`${data.backend_url}/api/v1/details`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({ url: formData.get('url') }),
		})
		playingSong = await response.json()
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
			<audio bind:this={audio} src={playingSong.play_url} controls autoplay></audio>
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

	<!-- <Card.Root>
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
	  </Card.Root> -->

</div>
