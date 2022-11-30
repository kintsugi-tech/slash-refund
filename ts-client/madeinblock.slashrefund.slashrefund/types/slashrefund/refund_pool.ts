/* eslint-disable */
import { Coin } from "../cosmos/base/v1beta1/coin";
import { Writer, Reader } from "protobufjs/minimal";

export const protobufPackage = "madeinblock.slashrefund.slashrefund";

/** TODO: to account for more than one token, Tokens and Shares must be a struct. */
export interface RefundPool {
  operatorAddress: string;
  tokens: Coin | undefined;
  shares: string;
}

const baseRefundPool: object = { operatorAddress: "", shares: "" };

export const RefundPool = {
  encode(message: RefundPool, writer: Writer = Writer.create()): Writer {
    if (message.operatorAddress !== "") {
      writer.uint32(10).string(message.operatorAddress);
    }
    if (message.tokens !== undefined) {
      Coin.encode(message.tokens, writer.uint32(18).fork()).ldelim();
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
          message.tokens = Coin.decode(reader, reader.uint32());
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
      message.tokens = Coin.fromJSON(object.tokens);
    } else {
      message.tokens = undefined;
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
    message.tokens !== undefined &&
      (obj.tokens = message.tokens ? Coin.toJSON(message.tokens) : undefined);
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
      message.tokens = Coin.fromPartial(object.tokens);
    } else {
      message.tokens = undefined;
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
