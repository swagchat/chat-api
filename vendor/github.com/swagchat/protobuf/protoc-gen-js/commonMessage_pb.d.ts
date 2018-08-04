// package: swagchat.protobuf
// file: commonMessage.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "./gogoproto/gogo_pb";

export class RoomIds extends jspb.Message {
  clearRoomIdsList(): void;
  getRoomIdsList(): Array<string>;
  setRoomIdsList(value: Array<string>): void;
  addRoomIds(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RoomIds.AsObject;
  static toObject(includeInstance: boolean, msg: RoomIds): RoomIds.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RoomIds, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RoomIds;
  static deserializeBinaryFromReader(message: RoomIds, reader: jspb.BinaryReader): RoomIds;
}

export namespace RoomIds {
  export type AsObject = {
    roomIdsList: Array<string>,
  }
}

export class UserIds extends jspb.Message {
  clearUserIdsList(): void;
  getUserIdsList(): Array<string>;
  setUserIdsList(value: Array<string>): void;
  addUserIds(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserIds.AsObject;
  static toObject(includeInstance: boolean, msg: UserIds): UserIds.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UserIds, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserIds;
  static deserializeBinaryFromReader(message: UserIds, reader: jspb.BinaryReader): UserIds;
}

export namespace UserIds {
  export type AsObject = {
    userIdsList: Array<string>,
  }
}

export class RoleIds extends jspb.Message {
  clearRoleIdsList(): void;
  getRoleIdsList(): Array<number>;
  setRoleIdsList(value: Array<number>): void;
  addRoleIds(value: number, index?: number): number;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RoleIds.AsObject;
  static toObject(includeInstance: boolean, msg: RoleIds): RoleIds.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RoleIds, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RoleIds;
  static deserializeBinaryFromReader(message: RoleIds, reader: jspb.BinaryReader): RoleIds;
}

export namespace RoleIds {
  export type AsObject = {
    roleIdsList: Array<number>,
  }
}

export class ErrorResponse extends jspb.Message {
  getMessage(): string;
  setMessage(value: string): void;

  getDeveloperMessage(): string;
  setDeveloperMessage(value: string): void;

  getInfo(): string;
  setInfo(value: string): void;

  clearInvalidParamsList(): void;
  getInvalidParamsList(): Array<InvalidParam>;
  setInvalidParamsList(value: Array<InvalidParam>): void;
  addInvalidParams(value?: InvalidParam, index?: number): InvalidParam;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ErrorResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ErrorResponse): ErrorResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ErrorResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ErrorResponse;
  static deserializeBinaryFromReader(message: ErrorResponse, reader: jspb.BinaryReader): ErrorResponse;
}

export namespace ErrorResponse {
  export type AsObject = {
    message: string,
    developerMessage: string,
    info: string,
    invalidParamsList: Array<InvalidParam.AsObject>,
  }
}

export class InvalidParam extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getReason(): string;
  setReason(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InvalidParam.AsObject;
  static toObject(includeInstance: boolean, msg: InvalidParam): InvalidParam.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: InvalidParam, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InvalidParam;
  static deserializeBinaryFromReader(message: InvalidParam, reader: jspb.BinaryReader): InvalidParam;
}

export namespace InvalidParam {
  export type AsObject = {
    name: string,
    reason: string,
  }
}

export class OrderInfo extends jspb.Message {
  getField(): string;
  setField(value: string): void;

  getOrder(): Order;
  setOrder(value: Order): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrderInfo.AsObject;
  static toObject(includeInstance: boolean, msg: OrderInfo): OrderInfo.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: OrderInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrderInfo;
  static deserializeBinaryFromReader(message: OrderInfo, reader: jspb.BinaryReader): OrderInfo;
}

export namespace OrderInfo {
  export type AsObject = {
    field: string,
    order: Order,
  }
}

export enum Order {
  ASC = 0,
  DESC = 1,
}

