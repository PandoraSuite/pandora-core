# Pandora Core

## :bug: Debugging

Debugging allows you to inspect the flow of the application, pause execution at breakpoints, and view values in real time. This can greatly improve development efficiency and help identify issues quickly.

To enable an efficient debugging experience while using `air`, `delve`, and optionally integrating with VSCode, follow the instructions below.

### :brain: Requirements

Make sure you have the following tools installed:

```sh
go install github.com/go-delve/delve/cmd/dlv@latest
go install github.com/air-verse/air@latest
```

### :gear: Setting up Debug with Air

Pandora-core includes an `air` configuration template that supports debugging out of the box.

> :page_facing_up: **Template file**: `.air.template.toml`

To configure it:

1. Copy the template:
   ```sh
   cp .air.toml.template .air.toml
   ```

2. Edit the `.air.toml` file:
   * **Remove** the `pre_cmd` line if you want to run the gRPC server.
   * **Adjust** the `cmd` field to point to the appropriate `main.go` entrypoint depending on which service you want to run. Example: `./cmd/http/main.go`

3. Run the application using:
   ```sh
   air
   ```

### :computer: VSCode

To debug using Visual Studio Code:

1. Open VSCode and go to the **Run and Debug** tab (Ctrl + Shift + D).
2. Click on **"Create a launch.json file"** and select **Go** as the environment.
3. Replace the generated content with the following configuration:

   ```json
   {
     "version": "0.2.0",
     "configurations": [
       {
         "name": "Connect to server",
         "type": "go",
         "request": "attach",
         "mode": "remote",
         "remotePath": "${workspaceFolder}",
         "port": 2345,
         "host": "127.0.0.1",
         "apiVersion": 2,
         "showLog": true,
         "trace": "verbose"
       }
     ]
   }
   ```

4. With `air` running the application, press `F5` in VSCode to start the debugging session.

### :wrench: Debugging in Other IDEs

If you're using another IDE (like Goland, Sublime Merge, Vim, etc.), feel free to contribute your own debug setup by documenting the necessary steps here. Keep in mind:

* Keep the configuration scoped to your IDE.
* Avoid modify project-specific files.
