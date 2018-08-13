// package: swagchat.protobuf
// file: messageMessage.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "./gogoproto/gogo_pb";

export class Messages extends jspb.Message {
  clearMessagesList(): void;
  getMessagesList(): Array<Message>;
  setMessagesList(value: Array<Message>): void;
  addMessages(value?: Message, index?: number): Message;

  hasAllcount(): boolean;
  clearAllcount(): void;
  getAllcount(): number | undefined;
  setAllcount(value: number): void;

  hasLimit(): boolean;
  clearLimit(): void;
  getLimit(): number | undefined;
  setLimit(value: number): void;

  hasOffset(): boolean;
  clearOffset(): void;
  getOffset(): number | undefined;
  setOffset(value: number): void;

  hasOrder(): boolean;
  clearOrder(): void;
  getOrder(): string | undefined;
  setOrder(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Messages.AsObject;
  static toObject(includeInstance: boolean, msg: Messages): Messages.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Messages, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Messages;
  static deserializeBinaryFromReader(message: Messages, reader: jspb.BinaryReader): Messages;
}

export namespace Messages {
  export type AsObject = {
    messagesList: Array<Message.AsObject>,
    allcount?: number,
    limit?: number,
    offset?: number,
    order?: string,
  }
}

export class Message extends jspb.Message {
  hasId(): boolean;
  clearId(): void;
  getId(): number | undefined;
  setId(value: number): void;

  hasMessageId(): boolean;
  clearMessageId(): void;
  getMessageId(): string | undefined;
  setMessageId(value: string): void;

  hasRoomId(): boolean;
  clearRoomId(): void;
  getRoomId(): string | undefined;
  setRoomId(value: string): void;

  hasUserId(): boolean;
  clearUserId(): void;
  getUserId(): string | undefined;
  setUserId(value: string): void;

  hasType(): boolean;
  clearType(): void;
  getType(): string | undefined;
  setType(value: string): void;

  hasPayload(): boolean;
  clearPayload(): void;
  getPayload(): Uint8Array | string;
  getPayload_asU8(): Uint8Array;
  getPayload_asB64(): string;
  setPayload(value: Uint8Array | string): void;

  hasRole(): boolean;
  clearRole(): void;
  getRole(): number | undefined;
  setRole(value: number): void;

  hasCreated(): boolean;
  clearCreated(): void;
  getCreated(): number | undefined;
  setCreated(value: number): void;

  hasModified(): boolean;
  clearModified(): void;
  getModified(): number | undefined;
  setModified(value: number): void;

  hasDeleted(): boolean;
  clearDeleted(): void;
  getDeleted(): number | undefined;
  setDeleted(value: number): void;

  hasEventName(): boolean;
  clearEventName(): void;
  getEventName(): string | undefined;
  setEventName(value: string): void;

  clearUserIdsList(): void;
  getUserIdsList(): Array<string>;
  setUserIdsList(value: Array<string>): void;
  addUserIds(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Message.AsObject;
  static toObject(includeInstance: boolean, msg: Message): Message.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Message, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Message;
  static deserializeBinaryFromReader(message: Message, reader: jspb.BinaryReader): Message;
}

export namespace Message {
  export type AsObject = {
    id?: number,
    messageId?: string,
    roomId?: string,
    userId?: string,
    type?: string,
    payload: Uint8Array | string,
    role?: number,
    created?: number,
    modified?: number,
    deleted?: number,
    eventName?: string,
    userIdsList: Array<string>,
  }
}

export class MessagePayload extends jspb.Message {
  hasText(): boolean;
  clearText(): void;
  getText(): string | undefined;
  setText(value: string): void;

  hasMime(): boolean;
  clearMime(): void;
  getMime(): string | undefined;
  setMime(value: string): void;

  hasFilename(): boolean;
  clearFilename(): void;
  getFilename(): string | undefined;
  setFilename(value: string): void;

  hasSourceurl(): boolean;
  clearSourceurl(): void;
  getSourceurl(): string | undefined;
  setSourceurl(value: string): void;

  hasThumbnailurl(): boolean;
  clearThumbnailurl(): void;
  getThumbnailurl(): string | undefined;
  setThumbnailurl(value: string): void;

  hasWidth(): boolean;
  clearWidth(): void;
  getWidth(): number | undefined;
  setWidth(value: number): void;

  hasHeight(): boolean;
  clearHeight(): void;
  getHeight(): number | undefined;
  setHeight(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MessagePayload.AsObject;
  static toObject(includeInstance: boolean, msg: MessagePayload): MessagePayload.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MessagePayload, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MessagePayload;
  static deserializeBinaryFromReader(message: MessagePayload, reader: jspb.BinaryReader): MessagePayload;
}

export namespace MessagePayload {
  export type AsObject = {
    text?: string,
    mime?: string,
    filename?: string,
    sourceurl?: string,
    thumbnailurl?: string,
    width?: number,
    height?: number,
  }
}

export class SendMessageRequest extends jspb.Message {
  hasMessageId(): boolean;
  clearMessageId(): void;
  getMessageId(): string | undefined;
  setMessageId(value: string): void;

  hasRoomId(): boolean;
  clearRoomId(): void;
  getRoomId(): string | undefined;
  setRoomId(value: string): void;

  hasUserId(): boolean;
  clearUserId(): void;
  getUserId(): string | undefined;
  setUserId(value: string): void;

  hasType(): boolean;
  clearType(): void;
  getType(): string | undefined;
  setType(value: string): void;

  hasPayload(): boolean;
  clearPayload(): void;
  getPayload(): Uint8Array | string;
  getPayload_asU8(): Uint8Array;
  getPayload_asB64(): string;
  setPayload(value: Uint8Array | string): void;

  hasRole(): boolean;
  clearRole(): void;
  getRole(): number | undefined;
  setRole(value: number): void;

  hasEventName(): boolean;
  clearEventName(): void;
  getEventName(): string | undefined;
  setEventName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SendMessageRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SendMessageRequest): SendMessageRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SendMessageRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SendMessageRequest;
  static deserializeBinaryFromReader(message: SendMessageRequest, reader: jspb.BinaryReader): SendMessageRequest;
}

export namespace SendMessageRequest {
  export type AsObject = {
    messageId?: string,
    roomId?: string,
    userId?: string,
    type?: string,
    payload: Uint8Array | string,
    role?: number,
    eventName?: string,
  }
}

