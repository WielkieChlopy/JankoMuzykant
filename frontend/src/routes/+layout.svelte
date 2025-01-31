<script lang="ts">
	import '../app.css';
	import { i18n } from '$lib/i18n';
	import { ParaglideJS } from '@inlang/paraglide-sveltekit';
	import { Toaster } from "$lib/components/ui/sonner/index.js";

	let { children } = $props();
</script>

<Toaster richColors closeButton/>

<ParaglideJS {i18n}>
	{@render children()}
</ParaglideJS>
