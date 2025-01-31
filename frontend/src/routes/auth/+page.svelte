<script lang="ts">
	import { Button } from "$lib/components/ui/button/index.js";
	import { Icons } from "$lib/components/icons/icons.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import { Label } from "$lib/components/ui/label/index.js";
	import { cn } from "$lib/utils.js";
	import { enhance } from '$app/forms';
	import * as m from '$lib/paraglide/messages.js';
	import { toast } from "svelte-sonner";

	let className: string | undefined | null = $state(undefined);
	export { className as class };

	let isLoading = $state(false)
	
	let { form } = $props();

	$effect(() => {
		if (form?.error) {
			toast.error(form.error);
		}
	});
</script>

<div
	class="container relative hidden h-[800px] flex-col items-center justify-center md:grid lg:max-w-none  lg:px-0"
>
	<div class="lg:p-8">
		<div class="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
			<div class="flex flex-col space-y-2 text-center">
				<h1 class="text-2xl font-semibold tracking-tight">{m.login()}</h1>
				<p class="text-muted-foreground text-sm">
					{m.enter_email()}
				</p>
			</div>
			
			<div class={cn("grid gap-6", className)}>
				<form 
					method="POST"
					action="?/login" 
					use:enhance={() => {
					isLoading = true;
			
						return async ({ update }) => {
							await update();
							isLoading = false;
						};
					}}
				>
					<div class="grid gap-2">
						<div class="grid gap-1">
							<Label class="sr-only" for="username">{m.username()}</Label>
							<Input
								id="username"
								name="username"
								placeholder="username"
								type="text"
								autocapitalize="none"
								autocorrect="off"
								disabled={isLoading}
								required
							/>
							<Label class="sr-only" for="password">{m.password()}</Label>
							<Input
								id="password"
								name="password"
								placeholder="********"
								type="password"
								disabled={isLoading}
								required
							/>
						</div>
						<Button type="submit" disabled={isLoading}>
							{#if isLoading}
								<Icons.spinner class="mr-2 h-4 w-4 animate-spin" />
							{/if}
							{m.login()}
						</Button>
					</div>
				</form>
			</div>
		</div>
	</div>
</div>