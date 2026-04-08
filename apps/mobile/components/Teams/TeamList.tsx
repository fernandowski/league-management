import {Animated, View} from "react-native";
import ScrollView = Animated.ScrollView;
import {Team} from "@/components/Teams/TeamOverview";
import { AppCard } from "@/components/ui/AppCard";
import { AppText } from "@/components/ui/AppText";

interface props {
    data: Team[]
}
export default function TeamList (props: props) {
    return (
        <View style={{flex: 1}}>
            <ScrollView>
                {
                    props.data.map((item) => (
                        <AppCard style={{marginBottom: 8, marginLeft: 1, marginRight: 1}} key={item.id}>
                            <AppCard.Content>
                                <AppText>{item.name}</AppText>
                            </AppCard.Content>
                        </AppCard>
                    ))
                }
            </ScrollView>
        </View>
    )
}
