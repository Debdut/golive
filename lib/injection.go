package lib

// Append script to HTML content
func reloadScript() string {
	return `
		<script>
			const source = new EventSource('/reload');
			source.onmessage = function(event) {
				if (event.data === 'reload') {
					window.location.reload();
				}
			};
		</script>
	`
}
