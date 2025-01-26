<script lang="ts">
	import * as Card from "@/components/ui/card/index.js";
	let { data } = $props();
    console.log(data)
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
			<audio id="player" src={data.url} controls autoplay></audio>
		</Card.Content>
	  </Card.Root>
</div>
