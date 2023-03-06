package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/fentec-project/gofe/abe"
	"github.com/spf13/cobra"
	"log"
	"sfs-go/internal/encrypt/pre/file"
	"sfs-go/internal/fabric/sdkInit"
	"sfs-go/internal/tools"
	"strings"
)

type IdAndAttr struct {
	Id    string
	Attrs []string
}

// cpabe represents the cpabe command
var cpabe = &cobra.Command{
	Use:   "abe",
	Short: "cp-abe",
	Long:  "cp-abe Related operations",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var add = &cobra.Command{
	Use:   "add",
	Short: "add",
	Long:  "add",
	Run: func(cmd *cobra.Command, args []string) {

		// 1. 检查MSP,并写入本地
		if mspBoolean != "" {
			_, err := abe.BooleanToMSP(mspBoolean, false)
			if err != nil {
				fmt.Printf("Failed to generate the policy: %v\n", err)
				return
			}
			// 存储
			file.WriteWithFile("config/msp", mspBoolean)
		}

		// id和属性集
		var idAttrs []IdAndAttr

		idAttr1 := IdAndAttr{
			Id:    id1,
			Attrs: attr1,
		}

		idAttr2 := IdAndAttr{
			Id:    id2,
			Attrs: attr2,
		}
		idAttr3 := IdAndAttr{
			Id:    id3,
			Attrs: attr3,
		}

		idAttrs = append(idAttrs, []IdAndAttr{idAttr1, idAttr2, idAttr3}...)
		address := []string{authAddr1, authAddr2, authAddr3}

		uploadIdAndAttrs(address[0], idAttrs[0])
		uploadIdAndAttrs(address[1], idAttrs[1])
		uploadIdAndAttrs(address[2], idAttrs[2])

		// auth的地址存储本地
		file.WriteWithFile("config/auth.Address", strings.Join(address, ","))
	},
}

func uploadIdAndAttrs(address string, attrs IdAndAttr) {
	service := sdkInit.GetInstance().InitFabric()
	attrsBytes, _ := json.Marshal(attrs)
	_, err := service.InsertAbeAttrsAndId(address, tools.ByteToString(attrsBytes))
	if err != nil {
		fmt.Println("uploadIdAndAttrs err ", err)
	}
}

var generateAuth = &cobra.Command{
	Use:   "generate",
	Short: "generate",
	Long:  "generate",
	Run: func(cmd *cobra.Command, args []string) {
		maabe := abe.NewMAABE()

		myAddress := tools.GetMyAddress()
		attrs := queryIDAndAttrs(myAddress)

		auth1, err := maabe.NewMAABEAuth(attrs.Id, attrs.Attrs)
		if err != nil {
			fmt.Printf("Failed generation authority %s: %v\n", attrs.Id, err)
		}
		// 序列化存储auth
		authByes, err := auth1.EncodeAuth()
		if err != nil {
			fmt.Println("Encode auth err ", err)
		}
		file.WriteWithFile(fmt.Sprintf("config/%s", attrs.Id), tools.ByteToString(authByes))

		// 上传公钥
		epk, err := auth1.EncodePukeykeys()
		if err != nil {
			fmt.Println("encoe pukey keys error:", err)
		}
		err = uploadPubkey(myAddress, epk)
		if err != nil {
			fmt.Println("upload pubkey error ", err)
		}
	},
}

func uploadPubkey(address string, pkBytes []byte) error {
	service := sdkInit.GetInstance().InitFabric()
	_, err := service.InsertAbeAuthPK(address, pkBytes)
	if err != nil {
		fmt.Println("uploadPubkey err ", err)
	}
	return nil
}

func getAuthPubkeys(maabe *abe.MAABE, addr []string) []*abe.MAABEPubKey {
	service := sdkInit.GetInstance().InitFabric()

	var abePubKey []*abe.MAABEPubKey
	auth, _ := maabe.NewMAABEAuth("auth1", []string{"auth1-ct1", "auth1-ct2"})

	for i := 0; i < len(addr); i++ {
		authPKEncodeBytes, err := service.QueryAbeAuthPK(addr[i])
		if err != nil {
			fmt.Println("uploadPubkey err ", err)
			return nil
		}
		//authPKBytes := tools.StringToByte(string(authPKEncodeBytes))

		authPks, err := auth.DecodePubkeys(authPKEncodeBytes)
		if err != nil {
			fmt.Println("auth.DecodeAuth err ", err)
			return nil
		}

		abePubKey = append(abePubKey, &authPks)

	}
	return abePubKey
}

func queryIDAndAttrs(address string) IdAndAttr {

	service := sdkInit.GetInstance().InitFabric()
	attrsAndIdBytes, err := service.QueryAbeAttrsAndId(address)
	if err != nil {
		fmt.Println("queryIDAndAttrs err ", err)
	}

	var idAndAttr IdAndAttr
	err = json.Unmarshal(tools.StringToByte(string(attrsAndIdBytes)), &idAndAttr)
	if err != nil {
		fmt.Println("IdAndAttr querry Unmarshal err ", err)
	}
	return idAndAttr
}

var abeShare = &cobra.Command{
	Use:   "share",
	Short: "share",
	Long:  "share",
	Run: func(cmd *cobra.Command, args []string) {
		maabe := abe.NewMAABE()

		// 1 msp
		mspStr := getMSP()
		msp, err := abe.BooleanToMSP(mspStr, false)
		if err != nil {
			fmt.Printf("Failed to generate the policy: %v\n", err)
			return
		}
		// 2. 公钥集合
		addres := file.ReadWithFile("config/auth.Address")

		pks := getAuthPubkeys(maabe, strings.Split(addres, ","))

		// 3. 获取加密fd
		fd := getFilefd(shareABEFileID)
		fmt.Println("fd", fd)
		// 3. 执行加密操作
		ct, err := maabe.Encrypt(string(fd), msp, pks)
		if err != nil {
			fmt.Println("maabe.Encrypt err ", err)
			return
		}
		// 4. 上传摘要密文
		abeCipherBytes, err := maabe.EncodeABECipher(ct)
		if err != nil {
			fmt.Println("encode abe cipher err ", err)
			return
		}
		uploadABECt(strings.Split(addres, ","), shareABEFileID, tools.ByteToString(abeCipherBytes))
	},
}

func uploadABECt(addrs []string, id string, ct string) {

	service := sdkInit.GetInstance().InitFabric()

	for i := 0; i < len(addrs); i++ {
		_, err := service.InsertAbeShareAddressFile(addrs[i], id, ct)
		if err != nil {
			fmt.Println("uploadPubkey err ", err)
		}
	}
}

func getMSP() string {
	return file.ReadWithFile("config/msp")
}

func getFilefd(id string) []byte {
	_, fileEncryptEntity := QueryFile(id)
	fd, err := GetAESKey(fileEncryptEntity.FileEncryptCipher)
	if err != nil {
		log.Println("get aes encrypt key failure,", err)
		return nil
	}
	return fd
}

var (
	gid string
	// 权威中心id
	id1 string
	// 权威中心地址
	authAddr1 string
	// 属性集
	attr1 []string

	id2 string
	// 权威中心地址
	authAddr2 string
	// 属性集
	attr2 []string

	id3 string
	// 权威中心地址
	authAddr3 string
	// 属性集
	attr3 []string

	// MSP访问结构
	mspBoolean string

	// share file id
	shareABEFileID string
)

func init() {
	rootCmd.AddCommand(cpabe)

	cpabe.AddCommand(add)

	add.Flags().StringVarP(&gid, "gid", "", "", "authority's gid")

	add.Flags().StringVarP(&id1, "id1", "", "", "authority's address")
	add.Flags().StringVarP(&authAddr1, "authAddr1", "", "", "authority's address")
	add.Flags().StringSliceVarP(&attr1, "attr1", "", []string{}, "attribute set")

	add.Flags().StringVarP(&id2, "id2", "", "", "authority's address")
	add.Flags().StringVarP(&authAddr2, "authAddr2", "", "", "authority's address")
	add.Flags().StringSliceVarP(&attr2, "attr2", "", []string{}, "attribute set")

	add.Flags().StringVarP(&id3, "id3", "", "", "authority's address")
	add.Flags().StringVarP(&authAddr3, "authAddr3", "", "", "authority's address")
	add.Flags().StringSliceVarP(&attr3, "attr3", "", []string{}, "attribute set")

	add.Flags().StringVarP(&mspBoolean, "msp", "m", "", "msp structure")

	cpabe.AddCommand(generateAuth)

	cpabe.AddCommand(abeShare)

	abeShare.Flags().StringVarP(&shareABEFileID, "fid", "f", "", "share file id")
	_ = abeShare.MarkFlagRequired("fid")

}
