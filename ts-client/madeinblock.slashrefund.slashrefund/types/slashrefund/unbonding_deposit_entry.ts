/* eslint-disable */
import { Writer, Reader } from "protobufjs/minimal";

export const protobufPackage = "madeinblock.slashrefund.slashrefund";

export interface UnbondingDepositEntry {
  creationHeight: number;
  completionTime: string;
  initialBalance: string;
  balance: string;
}

const baseUnbondingDepositEntry: object = {
  creationHeight: 0,
  completionTime: "",
  initialBalance: "",
  balance: "",
};

export const UnbondingDepositEntry = {
  encode(
    message: UnbondingDepositEntry,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.creationHeight !== 0) {
      writer.uint32(8).int32(message.creationHeight);
    }
    if (message.completionTime !== "") {
      writer.uint32(18).string(message.completionTime);
    }
    if (message.initialBalance !== "") {
      writer.uint32(26).string(message.initialBalance);
    }
    if (message.balance !== "") {
      writer.uint32(34).string(message.balance);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): UnbondingDepositEntry {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseUnbondingDepositEntry } as UnbondingDepositEntry;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creationHeight = reader.int32();
          break;
        case 2:
          message.completionTime = reader.string();
          break;
        case 3:
          message.initialBalance = reader.string();
          break;
        case 4:
          message.balance = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UnbondingDepositEntry {
    const message = { ...baseUnbondingDepositEntry } as UnbondingDepositEntry;
    if (object.creationHeight !== undefined && object.creationHeight !== null) {
      message.creationHeight = Number(object.creationHeight);
    } else {
      message.creationHeight = 0;
    }
    if (object.completionTime !== undefined && object.completionTime !== null) {
      message.completionTime = String(object.completionTime);
    } else {
      message.completionTime = "";
    }
    if (object.initialBalance !== undefined && object.initialBalance !== null) {
      message.initialBalance = String(object.initialBalance);
    } else {
      message.initialBalance = "";
    }
    if (object.balance !== undefined && object.balance !== null) {
      message.balance = String(object.balance);
    } else {
      message.balance = "";
    }
    return message;
  },

  toJSON(message: UnbondingDepositEntry): unknown {
    const obj: any = {};
    message.creationHeight !== undefined &&
      (obj.creationHeight = message.creationHeight);
    message.completionTime !== undefined &&
      (obj.completionTime = message.completionTime);
    message.initialBalance !== undefined &&
      (obj.initialBalance = message.initialBalance);
    message.balance !== undefined && (obj.balance = message.balance);
    return obj;
  },

  fromPartial(
    object: DeepPartial<UnbondingDepositEntry>
  ): UnbondingDepositEntry {
    const message = { ...baseUnbondingDepositEntry } as UnbondingDepositEntry;
    if (object.creationHeight !== undefined && object.creationHeight !== null) {
      message.creationHeight = object.creationHeight;
    } else {
      message.creationHeight = 0;
    }
    if (object.completionTime !== undefined && object.completionTime !== null) {
      message.completionTime = object.completionTime;
    } else {
      message.completionTime = "";
    }
    if (object.initialBalance !== undefined && object.initialBalance !== null) {
      message.initialBalance = object.initialBalance;
    } else {
      message.initialBalance = "";
    }
    if (object.balance !== undefined && object.balance !== null) {
      message.balance = object.balance;
    } else {
      message.balance = "";
    }
    return message;
  },
};

type Builtin = Date | Function | Uint8Array | string | number | undefined;
export type DeepPartial<T> = T extends Builtin
  ? T
  : T extends Array<infer U>
  ? Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U>
  ? ReadonlyArray<DeepPartial<U>>
  : T extends {}
  ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;
