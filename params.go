package main

type InitializeParams struct {
	/**
	 * The process Id of the parent process that started
	 * the server.
	 */
	ProcessID int32/*integer | null*/ `json:"processId"`
	/**
	 * Information about the client
	 *
	 * @since 3.15.0
	 */
	ClientInfo struct {
		/**
		 * The name of the client as defined by the client.
		 */
		Name string `json:"name"`
		/**
		 * The client's version as defined by the client.
		 */
		Version string `json:"version,omitempty"`
	} `json:"clientInfo,omitempty"`
	/**
	 * The locale the client is currently showing the user interface
	 * in. This must not necessarily be the locale of the operating
	 * system.
	 *
	 * Uses IETF language tags as the value's syntax
	 * (See https://en.wikipedia.org/wiki/IETF_language_tag)
	 *
	 * @since 3.16.0
	 */
	Locale string `json:"locale,omitempty"`
	/**
	 * The rootPath of the workspace. Is null
	 * if no folder is open.
	 *
	 * @deprecated in favour of rootUri.
	 */
	RootPath string/*string | null*/ `json:"rootPath,omitempty"`
	/**
	 * The rootUri of the workspace. Is null if no
	 * folder is open. If both `rootPath` and `rootUri` are set
	 * `rootUri` wins.
	 *
	 * @deprecated in favour of workspaceFolders.
	 */
	RootURI string/*DocumentUri | null*/ `json:"rootUri"`
	/**
	 * The capabilities provided by the client (editor or tool)
	 */
	// Capabilities ClientCapabilities `json:"capabilities"`
	/**
	 * User provided initialization options.
	 */
	InitializationOptions interface{} `json:"initializationOptions,omitempty"`
	/**
	 * The initial trace setting. If omitted trace is disabled ('off').
	 */
	Trace string/*'off' | 'messages' | 'verbose'*/ `json:"trace,omitempty"`
	/**
	 * The actual configured workspace folders.
	 */
	WorkspaceFolders []WorkspaceFolder/*WorkspaceFolder[] | null*/ `json:"workspaceFolders"`
}

type WorkspaceFolder struct {
	/**
	 * The associated URI for this workspace folder.
	 */
	URI string `json:"uri"`
	/**
	 * The name of the workspace folder. Used to refer to this
	 * workspace folder in the user interface.
	 */
	Name string `json:"name"`
}
