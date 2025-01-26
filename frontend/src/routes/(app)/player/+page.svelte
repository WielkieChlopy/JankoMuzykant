<script lang="ts">
	import * as Card from "@/components/ui/card/index.js";
	let { data } = $props();

	let playingSong = $state({
		url: '',
		play_url: '',
		title: '',
		duration: ''
	})

	let audio: HTMLAudioElement | null = null

	$effect(() => {
		if (data.queue.length > 0) {
			console.log('data.queue', data.queue)
			playingSong = data.queue[0]
			if (audio) {
				audio.src = playingSong.play_url
				audio.load()
			}
		}
	})
</script>

<div>
	<form method="POST" action="?/add_song">
		<input type="text" name="url" />
		<button type="submit">Add</button>
	</form>
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
	<Card.Root>
		<Card.Header>
		  <Card.Title>Player</Card.Title>
		  <Card.Description>Playing</Card.Description>
		</Card.Header>
		<Card.Content>
			<audio bind:this={audio} src={playingSong.play_url} controls autoplay></audio>
		</Card.Content>
	  </Card.Root>
</div>
