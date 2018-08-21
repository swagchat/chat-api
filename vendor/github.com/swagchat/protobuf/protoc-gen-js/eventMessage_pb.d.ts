// package: swagchat.protobuf
// file: eventMessage.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "./gogoproto/gogo_pb";
import * as roomMessage_pb from "./roomMessage_pb";

export class EventData extends jspb.Message {
  hasType(): boolean;
  clearType(): void;
  getType(): EventType | undefined;
  setType(value: EventType): void;

  hasData(): boolean;
  clearData(): void;
  getData(): Uint8Array | string;
  getData_asU8(): Uint8Array;
  getData_asB64(): string;
  setData(value: Uint8Array | string): void;

  clearUserIdsList(): void;
  getUserIdsList(): Array<string>;
  setUserIdsList(value: Array<string>): void;
  addUserIds(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EventData.AsObject;
  static toObject(includeInstance: boolean, msg: EventData): EventData.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EventData, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EventData;
  static deserializeBinaryFromReader(message: EventData, reader: jspb.BinaryReader): EventData;
}

export namespace EventData {
  export type AsObject = {
    type?: EventType,
    data: Uint8Array | string,
    userIdsList: Array<string>,
  }
}

export class UserJoinEventPayload extends jspb.Message {
  hasRoomId(): boolean;
  clearRoomId(): void;
  getRoomId(): string | undefined;
  setRoomId(value: string): void;

  clearUsersList(): void;
  getUsersList(): Array<roomMessage_pb.MiniUser>;
  setUsersList(value: Array<roomMessage_pb.MiniUser>): void;
  addUsers(value?: roomMessage_pb.MiniUser, index?: number): roomMessage_pb.MiniUser;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserJoinEventPayload.AsObject;
  static toObject(includeInstance: boolean, msg: UserJoinEventPayload): UserJoinEventPayload.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UserJoinEventPayload, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserJoinEventPayload;
  static deserializeBinaryFromReader(message: UserJoinEventPayload, reader: jspb.BinaryReader): UserJoinEventPayload;
}

export namespace UserJoinEventPayload {
  export type AsObject = {
    roomId?: string,
    usersList: Array<roomMessage_pb.MiniUser.AsObject>,
  }
}

export enum EventType {
  EMPTYEVENT = 0,
  MESSAGEEVENT = 1,
  USERJOINEVENT = 2,
}

