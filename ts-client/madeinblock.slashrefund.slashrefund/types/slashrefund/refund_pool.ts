/* eslint-disable */
import { Writer, Reader } from "protobufjs/minimal";

export const protobufPackage = "madeinblock.slashrefund.slashrefund";

export interface RefundPool {
  operatorAddress: string;
  tokens: string;
  shares: string;
}

const baseRefundPool: object = { operatorAddress: "", tokens: "", shares: "" };

export const RefundPool = {
  encode(message: RefundPool, writer: Writer = Writer.create()): Writer {
    if (message.operatorAddress !== "") {
      writer.uint32(10).string(message.operatorAddress);
    }
    if (message.tokens !== "") {
      writer.uint32(18).string(message.tokens);
    }
    if (message.shares !== "") {
      writer.uint32(26).string(message.shares);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): RefundPool {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseRefundPool } as RefundPool;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.operatorAddress = reader.string();
          break;
        case 2:
          message.tokens = reader.string();
          break;
        case 3:
          message.shares = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RefundPool {
    const message = { ...baseRefundPool } as RefundPool;
    if (
      object.operatorAddress !== undefined &&
      object.operatorAddress !== null
    ) {
      message.operatorAddress = String(object.operatorAddress);
    } else {
      message.operatorAddress = "";
    }
    if (object.tokens !== undefined && object.tokens !== null) {
      message.tokens = String(object.tokens);
    } else {
      message.tokens = "";
    }
    if (object.shares !== undefined && object.shares !== null) {
      message.shares = String(object.shares);
    } else {
      message.shares = "";
    }
    return message;
  },

  toJSON(message: RefundPool): unknown {
    const obj: any = {};
    message.operatorAddress !== undefined &&
      (obj.operatorAddress = message.operatorAddress);
    message.tokens !== undefined && (obj.tokens = message.tokens);
    message.shares !== undefined && (obj.shares = message.shares);
    return obj;
  },

  fromPartial(object: DeepPartial<RefundPool>): RefundPool {
    const message = { ...baseRefundPool } as RefundPool;
    if (
      object.operatorAddress !== undefined &&
      object.operatorAddress !== null
    ) {
      message.operatorAddress = object.operatorAddress;
    } else {
      message.operatorAddress = "";
    }
    if (object.tokens !== undefined && object.tokens !== null) {
      message.tokens = object.tokens;
    } else {
      message.tokens = "";
    }
    if (object.shares !== undefined && object.shares !== null) {
      message.shares = object.shares;
    } else {
      message.shares = "";
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
