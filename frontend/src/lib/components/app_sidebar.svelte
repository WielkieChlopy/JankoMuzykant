<script lang="ts">
    import * as Sidebar from "$lib/components/ui/sidebar/index.js";
</script>
    
<Sidebar.Root>
  <Sidebar.Header>
    <h1>Logo</h1>
  </Sidebar.Header>
    <Sidebar.Content>
      <Sidebar.Group>
        <Sidebar.GroupLabel>App</Sidebar.GroupLabel>
        <Sidebar.GroupContent>
          <Sidebar.Menu>
            <Sidebar.MenuItem>
              <Sidebar.MenuButton>
                Player
              </Sidebar.MenuButton>
              <Sidebar.MenuButton>
                Playlists
              </Sidebar.MenuButton>
              <Sidebar.MenuButton>
                Settings
              </Sidebar.MenuButton>
            </Sidebar.MenuItem>
          </Sidebar.Menu>
        </Sidebar.GroupContent>
      </Sidebar.Group>
    </Sidebar.Content>
  </Sidebar.Root>