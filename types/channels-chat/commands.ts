/**
 * Chat Channel Commands — `createChat`, `disposeChat`,
 * `addChatWorkspaceFolder`, and `removeChatWorkspaceFolder`.
 *
 * @module channels-chat/commands
 */

import type { URI } from '../common/state.js';
import type { BaseParams } from '../common/commands.js';
import type { Message } from './state.js';

// ─── createChat ──────────────────────────────────────────────────────────────

/**
 * Identifies a source chat and turn to fork from.
 */
export interface ChatForkSource {
  /** URI of the existing chat to fork from */
  chat: URI;
  /** Turn ID in the source chat; content up to and including this turn's response is copied */
  turnId: string;
}

/**
 * Creates a new chat within a session.
 *
 * @category Commands
 * @method createChat
 * @direction Client → Server
 * @messageType Request
 * @version 1
 */
export interface CreateChatParams extends BaseParams {
  /** Session URI containing the new chat. */
  channel: URI;
  /** Chat URI (client-chosen, e.g. `ahp-chat:/<uuid>`). */
  chat: URI;
  /** Optional initial message for the new chat. */
  initialMessage?: Message;
  /** Optional source chat and turn to fork from. */
  source?: ChatForkSource;
  /**
   * Initial working-directory subset for this chat. Every entry MUST be
   * present in the owning session's `workingDirectories`; the server MUST
   * reject any entry that is not. When absent, the chat inherits the full
   * session set. Forked chats (`source`) inherit the source chat's
   * `workingDirectories`; this field is ignored for forked chats.
   *
   * A client MUST NOT supply this field unless the agent advertises
   * {@link AgentCapabilities.multipleWorkspaceFolders}.
   */
  workingDirectories?: URI[];
}

// ─── disposeChat ─────────────────────────────────────────────────────────────

/**
 * Disposes a chat and cleans up server-side resources.
 *
 * @category Commands
 * @method disposeChat
 * @direction Client → Server
 * @messageType Request
 * @version 1
 */
export interface DisposeChatParams extends BaseParams {}

// ─── addChatWorkspaceFolder ──────────────────────────────────────────────────

/**
 * Grants this chat's agent tool access to a working directory, adding it to
 * the chat's {@link ChatState.workingDirectories | `workingDirectories`} subset.
 * The directory MUST already be present in the owning session's
 * `workingDirectories`; servers MUST reject with `InvalidParams` otherwise.
 *
 * Only valid when the agent advertises `multipleWorkspaceFolders`. Adding a
 * directory already in the chat's set is a no-op that still returns the current
 * full set.
 *
 * @category Commands
 * @method addChatWorkspaceFolder
 * @direction Client → Server
 * @messageType Request
 * @version 0.6.0
 */
export interface AddChatWorkspaceFolderParams extends BaseParams {
  /** Directory to grant tool access to. Must be in the session's `workingDirectories`. */
  folder: URI;
}

// ─── removeChatWorkspaceFolder ───────────────────────────────────────────────

/**
 * Revokes this chat's agent tool access to one of its working directories.
 * Analogous to the session-level `removeWorkspaceFolder`: the server
 * reconfigures the chat to the reduced subset and returns it. Removing a
 * directory not in the chat's set is a no-op that still returns the current
 * full set.
 *
 * Only valid when the agent advertises `multipleWorkspaceFolders`.
 *
 * @category Commands
 * @method removeChatWorkspaceFolder
 * @direction Client → Server
 * @messageType Request
 * @version 0.6.0
 */
export interface RemoveChatWorkspaceFolderParams extends BaseParams {
  /** Directory to revoke tool access to. */
  folder: URI;
}

/**
 * Result shared by `addChatWorkspaceFolder` and `removeChatWorkspaceFolder`:
 * the chat's full working-directory subset after the mutation.
 */
export interface ChatWorkspaceFolderResult {
  /** The chat's working directories after the mutation. */
  directories: URI[];
}

/**
 * Result of the `addChatWorkspaceFolder` command. See {@link ChatWorkspaceFolderResult}.
 */
export interface AddChatWorkspaceFolderResult extends ChatWorkspaceFolderResult {}

/**
 * Result of the `removeChatWorkspaceFolder` command. See {@link ChatWorkspaceFolderResult}.
 */
export interface RemoveChatWorkspaceFolderResult extends ChatWorkspaceFolderResult {}
