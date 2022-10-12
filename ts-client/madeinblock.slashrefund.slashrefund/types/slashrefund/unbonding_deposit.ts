/* eslint-disable */
import { Timestamp } from "../google/protobuf/timestamp";
import * as Long from "long";
import { util, configure, Writer, Reader } from "protobufjs/minimal";
import { Deposit } from "../slashrefund/deposit";

export const protobufPackage = "madeinblock.slashrefund.slashrefund";

export interface UnbondingDeposit {
  id: number;
  unbondingStart: Date | undefined;
  deposit: Deposit | undefined;
}

const baseUnbondingDeposit: object = { id: 0 };

export const UnbondingDeposit = {
  encode(message: UnbondingDeposit, writer: Writer = Writer.create()): Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint64(message.id);
    }
    if (message.unbondingStart !== undefined) {
      Timestamp.encode(
        toTimestamp(message.unbondingStart),
        writer.uint32(18).fork()
      ).ldelim();
    }
    if (message.deposit !== undefined) {
      Deposit.encode(message.deposit, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): UnbondingDeposit {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseUnbondingDeposit } as UnbondingDeposit;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.unbondingStart = fromTimestamp(
            Timestamp.decode(reader, reader.uint32())
          );
          break;
        case 3:
          message.deposit = Deposit.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UnbondingDeposit {
    const message = { ...baseUnbondingDeposit } as UnbondingDeposit;
    if (object.id !== undefined && object.id !== null) {
      message.id = Number(object.id);
    } else {
      message.id = 0;
    }
    if (object.unbondingStart !== undefined && object.unbondingStart !== null) {
      message.unbondingStart = fromJsonTimestamp(object.unbondingStart);
    } else {
      message.unbondingStart = undefined;
    }
    if (object.deposit !== undefined && object.deposit !== null) {
      message.deposit = Deposit.fromJSON(object.deposit);
    } else {
      message.deposit = undefined;
    }
    return message;
  },

  toJSON(message: UnbondingDeposit): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.unbondingStart !== undefined &&
      (obj.unbondingStart =
        message.unbondingStart !== undefined
          ? message.unbondingStart.toISOString()
          : null);
    message.deposit !== undefined &&
      (obj.deposit = message.deposit
        ? Deposit.toJSON(message.deposit)
        : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<UnbondingDeposit>): UnbondingDeposit {
    const message = { ...baseUnbondingDeposit } as UnbondingDeposit;
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    } else {
      message.id = 0;
    }
    if (object.unbondingStart !== undefined && object.unbondingStart !== null) {
      message.unbondingStart = object.unbondingStart;
    } else {
      message.unbondingStart = undefined;
    }
    if (object.deposit !== undefined && object.deposit !== null) {
      message.deposit = Deposit.fromPartial(object.deposit);
    } else {
      message.deposit = undefined;
    }
    return message;
  },
};

declare var self: any | undefined;
declare var window: any | undefined;
var globalThis: any = (() => {
  if (typeof globalThis !== "undefined") return globalThis;
  if (typeof self !== "undefined") return self;
  if (typeof window !== "undefined") return window;
  if (typeof global !== "undefined") return global;
  throw "Unable to locate global object";
})();

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

function toTimestamp(date: Date): Timestamp {
  const seconds = date.getTime() / 1_000;
  const nanos = (date.getTime() % 1_000) * 1_000_000;
  return { seconds, nanos };
}

function fromTimestamp(t: Timestamp): Date {
  let millis = t.seconds * 1_000;
  millis += t.nanos / 1_000_000;
  return new Date(millis);
}

function fromJsonTimestamp(o: any): Date {
  if (o instanceof Date) {
    return o;
  } else if (typeof o === "string") {
    return new Date(o);
  } else {
    return fromTimestamp(Timestamp.fromJSON(o));
  }
}

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (util.Long !== Long) {
  util.Long = Long as any;
  configure();
}
