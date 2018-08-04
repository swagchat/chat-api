// package: swagchat.protobuf
// file: userRoleMessage.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "./gogoproto/gogo_pb";

export class UserRole extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  getRole(): number;
  setRole(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserRole.AsObject;
  static toObject(includeInstance: boolean, msg: UserRole): UserRole.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UserRole, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserRole;
  static deserializeBinaryFromReader(message: UserRole, reader: jspb.BinaryReader): UserRole;
}

export namespace UserRole {
  export type AsObject = {
    userId: string,
    role: number,
  }
}

export class CreateUserRolesRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  clearRolesList(): void;
  getRolesList(): Array<number>;
  setRolesList(value: Array<number>): void;
  addRoles(value: number, index?: number): number;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateUserRolesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateUserRolesRequest): CreateUserRolesRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateUserRolesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateUserRolesRequest;
  static deserializeBinaryFromReader(message: CreateUserRolesRequest, reader: jspb.BinaryReader): CreateUserRolesRequest;
}

export namespace CreateUserRolesRequest {
  export type AsObject = {
    userId: string,
    rolesList: Array<number>,
  }
}

export class DeleteUserRolesRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  clearRolesList(): void;
  getRolesList(): Array<number>;
  setRolesList(value: Array<number>): void;
  addRoles(value: number, index?: number): number;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteUserRolesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteUserRolesRequest): DeleteUserRolesRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteUserRolesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteUserRolesRequest;
  static deserializeBinaryFromReader(message: DeleteUserRolesRequest, reader: jspb.BinaryReader): DeleteUserRolesRequest;
}

export namespace DeleteUserRolesRequest {
  export type AsObject = {
    userId: string,
    rolesList: Array<number>,
  }
}

