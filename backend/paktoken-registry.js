const { Registry } = require("@cosmjs/proto-signing");
const path = require("path");
const bankTx = require(path.join(process.cwd(), "node_modules/cosmjs-types/cosmos/bank/v1beta1/tx.js"));

function encodeVarint(value) {
  let n = BigInt(value);
  const bytes = [];

  while (n > 127n) {
    bytes.push(Number((n & 127n) | 128n));
    n >>= 7n;
  }

  bytes.push(Number(n));
  return Uint8Array.from(bytes);
}

function encodeString(fieldNumber, value) {
  const data = Buffer.from(value, "utf8");
  return Buffer.concat([
    Buffer.from([fieldNumber << 3 | 2]),
    Buffer.from(encodeVarint(data.length)),
    data,
  ]);
}

function encodeUint64(fieldNumber, value) {
  return Buffer.concat([
    Buffer.from([fieldNumber << 3]),
    Buffer.from(encodeVarint(value)),
  ]);
}

function decodeVarint(bytes, offset) {
  let result = 0n;
  let shift = 0n;

  while (offset.value < bytes.length) {
    const byte = bytes[offset.value++];
    result |= BigInt(byte & 0x7f) << shift;

    if ((byte & 0x80) === 0) {
      return result;
    }

    shift += 7n;
  }

  throw new Error("Invalid protobuf varint");
}

function decodeMessage(bytes) {
  const result = {};
  const offset = { value: 0 };

  while (offset.value < bytes.length) {
    const tag = Number(decodeVarint(bytes, offset));
    const fieldNumber = tag >> 3;
    const wireType = tag & 7;

    if (wireType === 2) {
      const length = Number(decodeVarint(bytes, offset));
      const end = offset.value + length;
      const value = Buffer.from(bytes.slice(offset.value, end)).toString("utf8");
      offset.value = end;

      if (fieldNumber === 1) {
        result.creator = value;
      }
    } else if (wireType === 0) {
      const value = decodeVarint(bytes, offset);

      if (fieldNumber === 2) {
        result.projectId = value.toString();
      } else if (fieldNumber === 3) {
        result.amount = value.toString();
      }
    } else {
      throw new Error(`Unsupported protobuf wire type: ${wireType}`);
    }
  }

  return result;
}

function createCreateProjectMessageType() {
  return {
    encode(message) {
      const parts = [];

      if (message.creator) parts.push(encodeString(1, message.creator));
      if (message.name) parts.push(encodeString(2, message.name));
      if (message.totalSupply) parts.push(encodeString(3, message.totalSupply));
      if (message.owner) parts.push(encodeString(4, message.owner));

      if (message.tokenPrice !== undefined && message.tokenPrice !== null) {
        parts.push(encodeUint64(5, message.tokenPrice));
      }

      if (message.category) {
        parts.push(encodeString(6, message.category));
      }

      return {
        finish() {
          return Buffer.concat(parts);
        },
      };
    },

    decode(bytes) {
      return decodeMessage(bytes);
    },

    fromPartial(object) {
       return {
         creator: object.creator || "",
         name: object.name || "",
         totalSupply:
            object.totalSupply !== undefined
            ? String(object.totalSupply)
            : "",
         owner: object.owner || "",
         tokenPrice:
            object.tokenPrice !== undefined
            ? String(object.tokenPrice)
            : "0",
         category: object.category || "",
       };
    },

    create(object) {
      return this.fromPartial(object || {});
    },
  };
}

function createMessageType() {
  return {
    encode(message) {
      const parts = [];

      if (message.creator) {
        parts.push(encodeString(1, message.creator));
      }

      if (message.projectId !== undefined && message.projectId !== null) {
        parts.push(encodeUint64(2, message.projectId));
      }

      if (message.amount !== undefined && message.amount !== null) {
        parts.push(encodeUint64(3, message.amount));
      }

      return {
        finish() {
          return Buffer.concat(parts);
        },
      };
    },

    decode(bytes) {
      return decodeMessage(bytes);
    },

    fromPartial(object) {
      return {
        creator: object.creator || "",
        projectId:
          object.projectId !== undefined
            ? String(object.projectId)
            : "0",
        amount:
          object.amount !== undefined
            ? String(object.amount)
            : "0",
      };
    },

    create(object) {
      return this.fromPartial(object || {});
    },
  };
}

function createSetReleasedSupplyType() {
  return {
    encode(message) {
      const parts = [];

      if (message.creator) {
        parts.push(encodeString(1, message.creator));
      }

      if (message.projectId !== undefined && message.projectId !== null) {
        parts.push(encodeUint64(2, message.projectId));
      }

      // releasedSupply is a string (can be an amount or a percentage like "2%")
      if (message.releasedSupply) {
        parts.push(encodeString(3, String(message.releasedSupply)));
      }

      return {
        finish() {
          return Buffer.concat(parts);
        },
      };
    },

    decode(bytes) {
      const result = {};
      const offset = { value: 0 };

      while (offset.value < bytes.length) {
        const tag = Number(decodeVarint(bytes, offset));
        const fieldNumber = tag >> 3;
        const wireType = tag & 7;

        if (wireType === 2) {
          const length = Number(decodeVarint(bytes, offset));
          const end = offset.value + length;
          const value = Buffer.from(bytes.slice(offset.value, end)).toString("utf8");
          offset.value = end;

          if (fieldNumber === 1) {
            result.creator = value;
          } else if (fieldNumber === 3) {
            result.releasedSupply = value;
          }
        } else if (wireType === 0) {
          const value = decodeVarint(bytes, offset);
          if (fieldNumber === 2) {
            result.projectId = value.toString();
          }
        }
      }

      return result;
    },

    fromPartial(object) {
      return {
        creator: object.creator || "",
        projectId:
          object.projectId !== undefined ? String(object.projectId) : "0",
        releasedSupply:
          object.releasedSupply !== undefined
            ? String(object.releasedSupply)
            : "0",
      };
    },

    create(object) {
      return this.fromPartial(object || {});
    },
  };
}

const MsgCreateProject = createCreateProjectMessageType();
const MsgBuyTokens = createMessageType();
const MsgDistributeRevenue = createMessageType();
const MsgSetReleasedSupply = createSetReleasedSupplyType();

const registry = new Registry([
  ["/cosmos.bank.v1beta1.MsgSend", bankTx.MsgSend],
  [
    "/paktoken.tokenization.v1.MsgCreateProject",
    MsgCreateProject,
  ],
  [
    "/paktoken.tokenization.v1.MsgBuyTokens",
    MsgBuyTokens,
  ],
  [
    "/paktoken.tokenization.v1.MsgDistributeRevenue",
    MsgDistributeRevenue,
  ],
  [
    "/paktoken.tokenization.v1.MsgSetReleasedSupply",
    MsgSetReleasedSupply,
  ],
]);

module.exports = { registry };
