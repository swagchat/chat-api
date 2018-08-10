// package: swagchat.protobuf
// file: deviceMessage.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "./gogoproto/gogo_pb";

export class Device extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  getPlatform(): Platform;
  setPlatform(value: Platform): void;

  getToken(): string;
  setToken(value: string): void;

  getNotificationDeviceId(): string;
  setNotificationDeviceId(value: string): void;

  getDeleted(): number;
  setDeleted(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Device.AsObject;
  static toObject(includeInstance: boolean, msg: Device): Device.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Device, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Device;
  static deserializeBinaryFromReader(message: Device, reader: jspb.BinaryReader): Device;
}

export namespace Device {
  export type AsObject = {
    userId: string,
    platform: Platform,
    token: string,
    notificationDeviceId: string,
    deleted: number,
  }
}

export class CreateDeviceRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  getPlatform(): Platform;
  setPlatform(value: Platform): void;

  getToken(): string;
  setToken(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateDeviceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateDeviceRequest): CreateDeviceRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateDeviceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateDeviceRequest;
  static deserializeBinaryFromReader(message: CreateDeviceRequest, reader: jspb.BinaryReader): CreateDeviceRequest;
}

export namespace CreateDeviceRequest {
  export type AsObject = {
    userId: string,
    platform: Platform,
    token: string,
  }
}

export class RetrieveDevicesRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RetrieveDevicesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RetrieveDevicesRequest): RetrieveDevicesRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RetrieveDevicesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RetrieveDevicesRequest;
  static deserializeBinaryFromReader(message: RetrieveDevicesRequest, reader: jspb.BinaryReader): RetrieveDevicesRequest;
}

export namespace RetrieveDevicesRequest {
  export type AsObject = {
    userId: string,
  }
}

export class DevicesResponse extends jspb.Message {
  clearDevicesList(): void;
  getDevicesList(): Array<Device>;
  setDevicesList(value: Array<Device>): void;
  addDevices(value?: Device, index?: number): Device;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DevicesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DevicesResponse): DevicesResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DevicesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DevicesResponse;
  static deserializeBinaryFromReader(message: DevicesResponse, reader: jspb.BinaryReader): DevicesResponse;
}

export namespace DevicesResponse {
  export type AsObject = {
    devicesList: Array<Device.AsObject>,
  }
}

export class DeleteDeviceRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  getPlatform(): Platform;
  setPlatform(value: Platform): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteDeviceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteDeviceRequest): DeleteDeviceRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteDeviceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteDeviceRequest;
  static deserializeBinaryFromReader(message: DeleteDeviceRequest, reader: jspb.BinaryReader): DeleteDeviceRequest;
}

export namespace DeleteDeviceRequest {
  export type AsObject = {
    userId: string,
    platform: Platform,
  }
}

export enum Platform {
  PLATFORMNONE = 0,
  PLATFORMIOS = 1,
  PLATFORMANDROID = 2,
}

