// package: swagchat.protobuf
// file: roomMessage.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "./gogoproto/gogo_pb";
import * as commonMessage_pb from "./commonMessage_pb";
import * as messageMessage_pb from "./messageMessage_pb";

export class Room extends jspb.Message {
  hasId(): boolean;
  clearId(): void;
  getId(): number | undefined;
  setId(value: number): void;

  hasRoomId(): boolean;
  clearRoomId(): void;
  getRoomId(): string | undefined;
  setRoomId(value: string): void;

  hasUserId(): boolean;
  clearUserId(): void;
  getUserId(): string | undefined;
  setUserId(value: string): void;

  hasName(): boolean;
  clearName(): void;
  getName(): string | undefined;
  setName(value: string): void;

  hasPictureUrl(): boolean;
  clearPictureUrl(): void;
  getPictureUrl(): string | undefined;
  setPictureUrl(value: string): void;

  hasInformationUrl(): boolean;
  clearInformationUrl(): void;
  getInformationUrl(): string | undefined;
  setInformationUrl(value: string): void;

  hasType(): boolean;
  clearType(): void;
  getType(): RoomType | undefined;
  setType(value: RoomType): void;

  hasCanLeft(): boolean;
  clearCanLeft(): void;
  getCanLeft(): boolean | undefined;
  setCanLeft(value: boolean): void;

  hasSpeechMode(): boolean;
  clearSpeechMode(): void;
  getSpeechMode(): SpeechMode | undefined;
  setSpeechMode(value: SpeechMode): void;

  hasMetaData(): boolean;
  clearMetaData(): void;
  getMetaData(): Uint8Array | string;
  getMetaData_asU8(): Uint8Array;
  getMetaData_asB64(): string;
  setMetaData(value: Uint8Array | string): void;

  hasAvailableMessageTypes(): boolean;
  clearAvailableMessageTypes(): void;
  getAvailableMessageTypes(): string | undefined;
  setAvailableMessageTypes(value: string): void;

  hasLastMessage(): boolean;
  clearLastMessage(): void;
  getLastMessage(): string | undefined;
  setLastMessage(value: string): void;

  hasLastMessageUpdated(): boolean;
  clearLastMessageUpdated(): void;
  getLastMessageUpdated(): number | undefined;
  setLastMessageUpdated(value: number): void;

  hasMessageCount(): boolean;
  clearMessageCount(): void;
  getMessageCount(): number | undefined;
  setMessageCount(value: number): void;

  hasNotificationTopicId(): boolean;
  clearNotificationTopicId(): void;
  getNotificationTopicId(): string | undefined;
  setNotificationTopicId(value: string): void;

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

  clearUsersList(): void;
  getUsersList(): Array<MiniUser>;
  setUsersList(value: Array<MiniUser>): void;
  addUsers(value?: MiniUser, index?: number): MiniUser;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Room.AsObject;
  static toObject(includeInstance: boolean, msg: Room): Room.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Room, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Room;
  static deserializeBinaryFromReader(message: Room, reader: jspb.BinaryReader): Room;
}

export namespace Room {
  export type AsObject = {
    id?: number,
    roomId?: string,
    userId?: string,
    name?: string,
    pictureUrl?: string,
    informationUrl?: string,
    type?: RoomType,
    canLeft?: boolean,
    speechMode?: SpeechMode,
    metaData: Uint8Array | string,
    availableMessageTypes?: string,
    lastMessage?: string,
    lastMessageUpdated?: number,
    messageCount?: number,
    notificationTopicId?: string,
    created?: number,
    modified?: number,
    deleted?: number,
    usersList: Array<MiniUser.AsObject>,
  }
}

export class MiniUser extends jspb.Message {
  hasUserId(): boolean;
  clearUserId(): void;
  getUserId(): string | undefined;
  setUserId(value: string): void;

  hasName(): boolean;
  clearName(): void;
  getName(): string | undefined;
  setName(value: string): void;

  hasPictureUrl(): boolean;
  clearPictureUrl(): void;
  getPictureUrl(): string | undefined;
  setPictureUrl(value: string): void;

  hasInformationUrl(): boolean;
  clearInformationUrl(): void;
  getInformationUrl(): string | undefined;
  setInformationUrl(value: string): void;

  hasMetaData(): boolean;
  clearMetaData(): void;
  getMetaData(): Uint8Array | string;
  getMetaData_asU8(): Uint8Array;
  getMetaData_asB64(): string;
  setMetaData(value: Uint8Array | string): void;

  hasCanBlock(): boolean;
  clearCanBlock(): void;
  getCanBlock(): boolean | undefined;
  setCanBlock(value: boolean): void;

  hasLastAccessed(): boolean;
  clearLastAccessed(): void;
  getLastAccessed(): number | undefined;
  setLastAccessed(value: number): void;

  hasRuDisplay(): boolean;
  clearRuDisplay(): void;
  getRuDisplay(): boolean | undefined;
  setRuDisplay(value: boolean): void;

  hasCreated(): boolean;
  clearCreated(): void;
  getCreated(): number | undefined;
  setCreated(value: number): void;

  hasModified(): boolean;
  clearModified(): void;
  getModified(): number | undefined;
  setModified(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MiniUser.AsObject;
  static toObject(includeInstance: boolean, msg: MiniUser): MiniUser.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MiniUser, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MiniUser;
  static deserializeBinaryFromReader(message: MiniUser, reader: jspb.BinaryReader): MiniUser;
}

export namespace MiniUser {
  export type AsObject = {
    userId?: string,
    name?: string,
    pictureUrl?: string,
    informationUrl?: string,
    metaData: Uint8Array | string,
    canBlock?: boolean,
    lastAccessed?: number,
    ruDisplay?: boolean,
    created?: number,
    modified?: number,
  }
}

export class CreateRoomRequest extends jspb.Message {
  hasRoomId(): boolean;
  clearRoomId(): void;
  getRoomId(): string | undefined;
  setRoomId(value: string): void;

  hasUserId(): boolean;
  clearUserId(): void;
  getUserId(): string | undefined;
  setUserId(value: string): void;

  hasName(): boolean;
  clearName(): void;
  getName(): string | undefined;
  setName(value: string): void;

  hasPictureUrl(): boolean;
  clearPictureUrl(): void;
  getPictureUrl(): string | undefined;
  setPictureUrl(value: string): void;

  hasInformationUrl(): boolean;
  clearInformationUrl(): void;
  getInformationUrl(): string | undefined;
  setInformationUrl(value: string): void;

  hasType(): boolean;
  clearType(): void;
  getType(): RoomType | undefined;
  setType(value: RoomType): void;

  hasCanLeft(): boolean;
  clearCanLeft(): void;
  getCanLeft(): boolean | undefined;
  setCanLeft(value: boolean): void;

  hasSpeechMode(): boolean;
  clearSpeechMode(): void;
  getSpeechMode(): SpeechMode | undefined;
  setSpeechMode(value: SpeechMode): void;

  hasMetaData(): boolean;
  clearMetaData(): void;
  getMetaData(): Uint8Array | string;
  getMetaData_asU8(): Uint8Array;
  getMetaData_asB64(): string;
  setMetaData(value: Uint8Array | string): void;

  hasAvailableMessageTypes(): boolean;
  clearAvailableMessageTypes(): void;
  getAvailableMessageTypes(): string | undefined;
  setAvailableMessageTypes(value: string): void;

  clearUserIdsList(): void;
  getUserIdsList(): Array<string>;
  setUserIdsList(value: Array<string>): void;
  addUserIds(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateRoomRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateRoomRequest): CreateRoomRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateRoomRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateRoomRequest;
  static deserializeBinaryFromReader(message: CreateRoomRequest, reader: jspb.BinaryReader): CreateRoomRequest;
}

export namespace CreateRoomRequest {
  export type AsObject = {
    roomId?: string,
    userId?: string,
    name?: string,
    pictureUrl?: string,
    informationUrl?: string,
    type?: RoomType,
    canLeft?: boolean,
    speechMode?: SpeechMode,
    metaData: Uint8Array | string,
    availableMessageTypes?: string,
    userIdsList: Array<string>,
  }
}

export class GetRoomsRequest extends jspb.Message {
  hasLimit(): boolean;
  clearLimit(): void;
  getLimit(): number | undefined;
  setLimit(value: number): void;

  hasOffset(): boolean;
  clearOffset(): void;
  getOffset(): number | undefined;
  setOffset(value: number): void;

  clearOrdersList(): void;
  getOrdersList(): Array<commonMessage_pb.OrderInfo>;
  setOrdersList(value: Array<commonMessage_pb.OrderInfo>): void;
  addOrders(value?: commonMessage_pb.OrderInfo, index?: number): commonMessage_pb.OrderInfo;

  hasUserId(): boolean;
  clearUserId(): void;
  getUserId(): string | undefined;
  setUserId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetRoomsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetRoomsRequest): GetRoomsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetRoomsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetRoomsRequest;
  static deserializeBinaryFromReader(message: GetRoomsRequest, reader: jspb.BinaryReader): GetRoomsRequest;
}

export namespace GetRoomsRequest {
  export type AsObject = {
    limit?: number,
    offset?: number,
    ordersList: Array<commonMessage_pb.OrderInfo.AsObject>,
    userId?: string,
  }
}

export class RoomsResponse extends jspb.Message {
  clearRoomsList(): void;
  getRoomsList(): Array<Room>;
  setRoomsList(value: Array<Room>): void;
  addRooms(value?: Room, index?: number): Room;

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

  clearOrdersList(): void;
  getOrdersList(): Array<commonMessage_pb.OrderInfo>;
  setOrdersList(value: Array<commonMessage_pb.OrderInfo>): void;
  addOrders(value?: commonMessage_pb.OrderInfo, index?: number): commonMessage_pb.OrderInfo;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RoomsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RoomsResponse): RoomsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RoomsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RoomsResponse;
  static deserializeBinaryFromReader(message: RoomsResponse, reader: jspb.BinaryReader): RoomsResponse;
}

export namespace RoomsResponse {
  export type AsObject = {
    roomsList: Array<Room.AsObject>,
    allcount?: number,
    limit?: number,
    offset?: number,
    ordersList: Array<commonMessage_pb.OrderInfo.AsObject>,
  }
}

export class GetRoomRequest extends jspb.Message {
  hasRoomId(): boolean;
  clearRoomId(): void;
  getRoomId(): string | undefined;
  setRoomId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetRoomRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetRoomRequest): GetRoomRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetRoomRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetRoomRequest;
  static deserializeBinaryFromReader(message: GetRoomRequest, reader: jspb.BinaryReader): GetRoomRequest;
}

export namespace GetRoomRequest {
  export type AsObject = {
    roomId?: string,
  }
}

export class UpdateRoomRequest extends jspb.Message {
  hasRoomId(): boolean;
  clearRoomId(): void;
  getRoomId(): string | undefined;
  setRoomId(value: string): void;

  hasName(): boolean;
  clearName(): void;
  getName(): string | undefined;
  setName(value: string): void;

  hasPictureUrl(): boolean;
  clearPictureUrl(): void;
  getPictureUrl(): string | undefined;
  setPictureUrl(value: string): void;

  hasInformationUrl(): boolean;
  clearInformationUrl(): void;
  getInformationUrl(): string | undefined;
  setInformationUrl(value: string): void;

  hasType(): boolean;
  clearType(): void;
  getType(): RoomType | undefined;
  setType(value: RoomType): void;

  hasCanLeft(): boolean;
  clearCanLeft(): void;
  getCanLeft(): boolean | undefined;
  setCanLeft(value: boolean): void;

  hasSpeechMode(): boolean;
  clearSpeechMode(): void;
  getSpeechMode(): SpeechMode | undefined;
  setSpeechMode(value: SpeechMode): void;

  hasMetaData(): boolean;
  clearMetaData(): void;
  getMetaData(): Uint8Array | string;
  getMetaData_asU8(): Uint8Array;
  getMetaData_asB64(): string;
  setMetaData(value: Uint8Array | string): void;

  hasAvailableMessageTypes(): boolean;
  clearAvailableMessageTypes(): void;
  getAvailableMessageTypes(): string | undefined;
  setAvailableMessageTypes(value: string): void;

  clearUserIdsList(): void;
  getUserIdsList(): Array<string>;
  setUserIdsList(value: Array<string>): void;
  addUserIds(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateRoomRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateRoomRequest): UpdateRoomRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateRoomRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateRoomRequest;
  static deserializeBinaryFromReader(message: UpdateRoomRequest, reader: jspb.BinaryReader): UpdateRoomRequest;
}

export namespace UpdateRoomRequest {
  export type AsObject = {
    roomId?: string,
    name?: string,
    pictureUrl?: string,
    informationUrl?: string,
    type?: RoomType,
    canLeft?: boolean,
    speechMode?: SpeechMode,
    metaData: Uint8Array | string,
    availableMessageTypes?: string,
    userIdsList: Array<string>,
  }
}

export class DeleteRoomRequest extends jspb.Message {
  hasRoomId(): boolean;
  clearRoomId(): void;
  getRoomId(): string | undefined;
  setRoomId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteRoomRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteRoomRequest): DeleteRoomRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteRoomRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteRoomRequest;
  static deserializeBinaryFromReader(message: DeleteRoomRequest, reader: jspb.BinaryReader): DeleteRoomRequest;
}

export namespace DeleteRoomRequest {
  export type AsObject = {
    roomId?: string,
  }
}

export class GetRoomMessagesRequest extends jspb.Message {
  hasLimit(): boolean;
  clearLimit(): void;
  getLimit(): number | undefined;
  setLimit(value: number): void;

  hasOffset(): boolean;
  clearOffset(): void;
  getOffset(): number | undefined;
  setOffset(value: number): void;

  clearOrdersList(): void;
  getOrdersList(): Array<commonMessage_pb.OrderInfo>;
  setOrdersList(value: Array<commonMessage_pb.OrderInfo>): void;
  addOrders(value?: commonMessage_pb.OrderInfo, index?: number): commonMessage_pb.OrderInfo;

  hasRoomId(): boolean;
  clearRoomId(): void;
  getRoomId(): string | undefined;
  setRoomId(value: string): void;

  clearRoleIdsList(): void;
  getRoleIdsList(): Array<number>;
  setRoleIdsList(value: Array<number>): void;
  addRoleIds(value: number, index?: number): number;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetRoomMessagesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetRoomMessagesRequest): GetRoomMessagesRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetRoomMessagesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetRoomMessagesRequest;
  static deserializeBinaryFromReader(message: GetRoomMessagesRequest, reader: jspb.BinaryReader): GetRoomMessagesRequest;
}

export namespace GetRoomMessagesRequest {
  export type AsObject = {
    limit?: number,
    offset?: number,
    ordersList: Array<commonMessage_pb.OrderInfo.AsObject>,
    roomId?: string,
    roleIdsList: Array<number>,
  }
}

export class RoomMessagesResponse extends jspb.Message {
  clearMessagesList(): void;
  getMessagesList(): Array<messageMessage_pb.Message>;
  setMessagesList(value: Array<messageMessage_pb.Message>): void;
  addMessages(value?: messageMessage_pb.Message, index?: number): messageMessage_pb.Message;

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

  clearOrdersList(): void;
  getOrdersList(): Array<commonMessage_pb.OrderInfo>;
  setOrdersList(value: Array<commonMessage_pb.OrderInfo>): void;
  addOrders(value?: commonMessage_pb.OrderInfo, index?: number): commonMessage_pb.OrderInfo;

  hasRoomId(): boolean;
  clearRoomId(): void;
  getRoomId(): string | undefined;
  setRoomId(value: string): void;

  clearRoleIdsList(): void;
  getRoleIdsList(): Array<number>;
  setRoleIdsList(value: Array<number>): void;
  addRoleIds(value: number, index?: number): number;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RoomMessagesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RoomMessagesResponse): RoomMessagesResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RoomMessagesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RoomMessagesResponse;
  static deserializeBinaryFromReader(message: RoomMessagesResponse, reader: jspb.BinaryReader): RoomMessagesResponse;
}

export namespace RoomMessagesResponse {
  export type AsObject = {
    messagesList: Array<messageMessage_pb.Message.AsObject>,
    allcount?: number,
    limit?: number,
    offset?: number,
    ordersList: Array<commonMessage_pb.OrderInfo.AsObject>,
    roomId?: string,
    roleIdsList: Array<number>,
  }
}

export enum RoomType {
  ONEONONEROOM = 0,
  PRIVATEROOM = 1,
  PUBLICROOM = 2,
  NOTICEROOM = 3,
}

export enum SpeechMode {
  SPEECHMODENONE = 0,
  WAKEUPWEBTOWEB = 1,
  WAKEUPWEBTOCLOUD = 2,
  WAKEUPCLOUDTOCLOUD = 3,
  ALWAYS = 4,
  MANUAL = 5,
}

