// package: swagchat.protobuf
// file: roomUserMessage.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "./gogoproto/gogo_pb";

export class RoomUser extends jspb.Message {
  hasRoomId(): boolean;
  clearRoomId(): void;
  getRoomId(): string | undefined;
  setRoomId(value: string): void;

  hasUserId(): boolean;
  clearUserId(): void;
  getUserId(): string | undefined;
  setUserId(value: string): void;

  hasUnreadCount(): boolean;
  clearUnreadCount(): void;
  getUnreadCount(): number | undefined;
  setUnreadCount(value: number): void;

  hasDisplay(): boolean;
  clearDisplay(): void;
  getDisplay(): boolean | undefined;
  setDisplay(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RoomUser.AsObject;
  static toObject(includeInstance: boolean, msg: RoomUser): RoomUser.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RoomUser, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RoomUser;
  static deserializeBinaryFromReader(message: RoomUser, reader: jspb.BinaryReader): RoomUser;
}

export namespace RoomUser {
  export type AsObject = {
    roomId?: string,
    userId?: string,
    unreadCount?: number,
    display?: boolean,
  }
}

export class CreateRoomUsersRequest extends jspb.Message {
  hasRoomId(): boolean;
  clearRoomId(): void;
  getRoomId(): string | undefined;
  setRoomId(value: string): void;

  clearUserIdsList(): void;
  getUserIdsList(): Array<string>;
  setUserIdsList(value: Array<string>): void;
  addUserIds(value: string, index?: number): string;

  hasDisplay(): boolean;
  clearDisplay(): void;
  getDisplay(): boolean | undefined;
  setDisplay(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateRoomUsersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateRoomUsersRequest): CreateRoomUsersRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateRoomUsersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateRoomUsersRequest;
  static deserializeBinaryFromReader(message: CreateRoomUsersRequest, reader: jspb.BinaryReader): CreateRoomUsersRequest;
}

export namespace CreateRoomUsersRequest {
  export type AsObject = {
    roomId?: string,
    userIdsList: Array<string>,
    display?: boolean,
  }
}

export class GetRoomUsersRequest extends jspb.Message {
  hasRoomId(): boolean;
  clearRoomId(): void;
  getRoomId(): string | undefined;
  setRoomId(value: string): void;

  clearUserIdsList(): void;
  getUserIdsList(): Array<string>;
  setUserIdsList(value: Array<string>): void;
  addUserIds(value: string, index?: number): string;

  clearRoleIdsList(): void;
  getRoleIdsList(): Array<number>;
  setRoleIdsList(value: Array<number>): void;
  addRoleIds(value: number, index?: number): number;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetRoomUsersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetRoomUsersRequest): GetRoomUsersRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetRoomUsersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetRoomUsersRequest;
  static deserializeBinaryFromReader(message: GetRoomUsersRequest, reader: jspb.BinaryReader): GetRoomUsersRequest;
}

export namespace GetRoomUsersRequest {
  export type AsObject = {
    roomId?: string,
    userIdsList: Array<string>,
    roleIdsList: Array<number>,
  }
}

export class RoomUsersResponse extends jspb.Message {
  clearUsersList(): void;
  getUsersList(): Array<RoomUser>;
  setUsersList(value: Array<RoomUser>): void;
  addUsers(value?: RoomUser, index?: number): RoomUser;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RoomUsersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RoomUsersResponse): RoomUsersResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RoomUsersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RoomUsersResponse;
  static deserializeBinaryFromReader(message: RoomUsersResponse, reader: jspb.BinaryReader): RoomUsersResponse;
}

export namespace RoomUsersResponse {
  export type AsObject = {
    usersList: Array<RoomUser.AsObject>,
  }
}

export class RoomUserIdsResponse extends jspb.Message {
  clearUserIdsList(): void;
  getUserIdsList(): Array<string>;
  setUserIdsList(value: Array<string>): void;
  addUserIds(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RoomUserIdsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RoomUserIdsResponse): RoomUserIdsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RoomUserIdsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RoomUserIdsResponse;
  static deserializeBinaryFromReader(message: RoomUserIdsResponse, reader: jspb.BinaryReader): RoomUserIdsResponse;
}

export namespace RoomUserIdsResponse {
  export type AsObject = {
    userIdsList: Array<string>,
  }
}

export class UpdateRoomUserRequest extends jspb.Message {
  hasRoomId(): boolean;
  clearRoomId(): void;
  getRoomId(): string | undefined;
  setRoomId(value: string): void;

  hasUserId(): boolean;
  clearUserId(): void;
  getUserId(): string | undefined;
  setUserId(value: string): void;

  hasUnreadCount(): boolean;
  clearUnreadCount(): void;
  getUnreadCount(): number | undefined;
  setUnreadCount(value: number): void;

  hasDisplay(): boolean;
  clearDisplay(): void;
  getDisplay(): boolean | undefined;
  setDisplay(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateRoomUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateRoomUserRequest): UpdateRoomUserRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateRoomUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateRoomUserRequest;
  static deserializeBinaryFromReader(message: UpdateRoomUserRequest, reader: jspb.BinaryReader): UpdateRoomUserRequest;
}

export namespace UpdateRoomUserRequest {
  export type AsObject = {
    roomId?: string,
    userId?: string,
    unreadCount?: number,
    display?: boolean,
  }
}

export class DeleteRoomUsersRequest extends jspb.Message {
  hasRoomId(): boolean;
  clearRoomId(): void;
  getRoomId(): string | undefined;
  setRoomId(value: string): void;

  clearUserIdsList(): void;
  getUserIdsList(): Array<string>;
  setUserIdsList(value: Array<string>): void;
  addUserIds(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteRoomUsersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteRoomUsersRequest): DeleteRoomUsersRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteRoomUsersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteRoomUsersRequest;
  static deserializeBinaryFromReader(message: DeleteRoomUsersRequest, reader: jspb.BinaryReader): DeleteRoomUsersRequest;
}

export namespace DeleteRoomUsersRequest {
  export type AsObject = {
    roomId?: string,
    userIdsList: Array<string>,
  }
}

