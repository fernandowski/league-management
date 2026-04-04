import {Button, Text, View} from "react-native";

export default function Index() {
  return (
    <View
      style={{
        flex: 1,
        justifyContent: "center",
        alignItems: "center",
      }}
    >
      <Text>I have edited </Text>
        <Button title={'click me'}></Button>
    </View>
  );
}
