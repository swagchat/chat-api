// package: swagchat.protobuf
// file: blockUserMessage.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "./gogoproto/gogo_pb";
import * as roomMessage_pb from "./roomMessage_pb";

export class BlockUser extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  getBlockUserId(): string;
  setBlockUserId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BlockUser.AsObject;
  static toObject(includeInstance: boolean, msg: BlockUser): BlockUser.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: BlockUser, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BlockUser;
  static deserializeBinaryFromReader(message: BlockUser, reader: jspb.BinaryReader): BlockUser;
}

export namespace BlockUser {
  export type AsObject = {
    userId: string,
    blockUserId: string,
  }
}

export class CreateBlockUsersRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  clearBlockUserIdsList(): void;
  getBlockUserIdsList(): Array<string>;
  setBlockUserIdsList(value: Array<string>): void;
  addBlockUserIds(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateBlockUsersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateBlockUsersRequest): CreateBlockUsersRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateBlockUsersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateBlockUsersRequest;
  static deserializeBinaryFromReader(message: CreateBlockUsersRequest, reader: jspb.BinaryReader): CreateBlockUsersRequest;
}

export namespace CreateBlockUsersRequest {
  export type AsObject = {
    userId: string,
    blockUserIdsList: Array<string>,
  }
}

export class GetBlockUsersRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetBlockUsersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetBlockUsersRequest): GetBlockUsersRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetBlockUsersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetBlockUsersRequest;
  static deserializeBinaryFromReader(message: GetBlockUsersRequest, reader: jspb.BinaryReader): GetBlockUsersRequest;
}

export namespace GetBlockUsersRequest {
  export type AsObject = {
    userId: string,
  }
}

export class BlockUsersResponse extends jspb.Message {
  clearBlockUsersList(): void;
  getBlockUsersList(): Array<roomMessage_pb.MiniUser>;
  setBlockUsersList(value: Array<roomMessage_pb.MiniUser>): void;
  addBlockUsers(value?: roomMessage_pb.MiniUser, index?: number): roomMessage_pb.MiniUser;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BlockUsersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: BlockUsersResponse): BlockUsersResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: BlockUsersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BlockUsersResponse;
  static deserializeBinaryFromReader(message: BlockUsersResponse, reader: jspb.BinaryReader): BlockUsersResponse;
}

export namespace BlockUsersResponse {
  export type AsObject = {
    blockUsersList: Array<roomMessage_pb.MiniUser.AsObject>,
  }
}

export class BlockUserIdsResponse extends jspb.Message {
  clearBlockUserIdsList(): void;
  getBlockUserIdsList(): Array<string>;
  setBlockUserIdsList(value: Array<string>): void;
  addBlockUserIds(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BlockUserIdsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: BlockUserIdsResponse): BlockUserIdsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: BlockUserIdsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BlockUserIdsResponse;
  static deserializeBinaryFromReader(message: BlockUserIdsResponse, reader: jspb.BinaryReader): BlockUserIdsResponse;
}

export namespace BlockUserIdsResponse {
  export type AsObject = {
    blockUserIdsList: Array<string>,
  }
}

export class GetBlockedUsersRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetBlockedUsersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetBlockedUsersRequest): GetBlockedUsersRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetBlockedUsersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetBlockedUsersRequest;
  static deserializeBinaryFromReader(message: GetBlockedUsersRequest, reader: jspb.BinaryReader): GetBlockedUsersRequest;
}

export namespace GetBlockedUsersRequest {
  export type AsObject = {
    userId: string,
  }
}

export class BlockedUsersResponse extends jspb.Message {
  clearBlockedUsersList(): void;
  getBlockedUsersList(): Array<roomMessage_pb.MiniUser>;
  setBlockedUsersList(value: Array<roomMessage_pb.MiniUser>): void;
  addBlockedUsers(value?: roomMessage_pb.MiniUser, index?: number): roomMessage_pb.MiniUser;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BlockedUsersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: BlockedUsersResponse): BlockedUsersResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: BlockedUsersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BlockedUsersResponse;
  static deserializeBinaryFromReader(message: BlockedUsersResponse, reader: jspb.BinaryReader): BlockedUsersResponse;
}

export namespace BlockedUsersResponse {
  export type AsObject = {
    blockedUsersList: Array<roomMessage_pb.MiniUser.AsObject>,
  }
}

export class BlockedUserIdsResponse extends jspb.Message {
  clearBlockedUserIdsList(): void;
  getBlockedUserIdsList(): Array<string>;
  setBlockedUserIdsList(value: Array<string>): void;
  addBlockedUserIds(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BlockedUserIdsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: BlockedUserIdsResponse): BlockedUserIdsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: BlockedUserIdsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BlockedUserIdsResponse;
  static deserializeBinaryFromReader(message: BlockedUserIdsResponse, reader: jspb.BinaryReader): BlockedUserIdsResponse;
}

export namespace BlockedUserIdsResponse {
  export type AsObject = {
    blockedUserIdsList: Array<string>,
  }
}

export class DeleteBlockUsersRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  clearBlockUserIdsList(): void;
  getBlockUserIdsList(): Array<string>;
  setBlockUserIdsList(value: Array<string>): void;
  addBlockUserIds(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteBlockUsersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteBlockUsersRequest): DeleteBlockUsersRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteBlockUsersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteBlockUsersRequest;
  static deserializeBinaryFromReader(message: DeleteBlockUsersRequest, reader: jspb.BinaryReader): DeleteBlockUsersRequest;
}

export namespace DeleteBlockUsersRequest {
  export type AsObject = {
    userId: string,
    blockUserIdsList: Array<string>,
  }
}

