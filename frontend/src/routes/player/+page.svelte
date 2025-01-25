<script lang="ts">
	let { data } = $props();
    console.log(data)
</script>

<div>
	<h1>Player</h1>
	<audio id="player" src={data.url} controls autoplay></audio>
</div>
